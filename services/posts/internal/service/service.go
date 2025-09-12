package service

import (
	"context"
	"encoding/json"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"soa-socialnetwork/services/common/backjob"
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
	"soa-socialnetwork/services/posts/internal/repos"
	"soa-socialnetwork/services/posts/internal/service/interceptors"
	"soa-socialnetwork/services/posts/internal/storage/postgres"
	pb "soa-socialnetwork/services/posts/proto"
	statsModels "soa-socialnetwork/services/stats/pkg/models"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsService struct {
	pb.UnimplementedPostsServiceServer

	Db          repos.Database
	JwtVerifier soajwt.Verifier

	outboxJob backjob.TickerJob
}

func New(cfg PostsServiceConfig) (PostsService, error) {
	db, err := postgres.NewPoolDatabase(postgres.ConnectionConfig{
		Host:     cfg.DbHost,
		User:     cfg.DbUser,
		Password: cfg.DbPassword,
		PoolSize: cfg.DbPoolSize,
	})

	if err != nil {
		return PostsService{}, err
	}

	return PostsService{
		Db:          &db,
		JwtVerifier: soajwt.NewVerifier(cfg.JwtPublicKey),
	}, nil
}

func (s *PostsService) Start() {
	outboxJobCallback := newCheckOutboxCallback(s.Db, 100)
	s.outboxJob = backjob.NewTickerJob(3*time.Second, outboxJobCallback)
	s.outboxJob.Run()
}

func (s *PostsService) EditPageSettings(ctx context.Context, req *pb.EditPageSettingsRequest) (*pb.Empty, error) {
	accountId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if accountId == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	if dom.AccountId(req.AccountId) != accountId.(dom.AccountId) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	pageData, err := conn.Pages().GetByAccountId(dom.AccountId(req.AccountId))

	if err != nil {
		return nil, err
	}

	err = conn.Pages().Edit(pageData.Id, repos.EditedPageSettings{
		VisibleForUnauthorized: opt.FromPointer(req.VisibleForUnauthorized),
		CommentsEnabled:        opt.FromPointer(req.CommentsEnabled),
		AnyoneCanPost:          opt.FromPointer(req.AnyoneCanPost),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *PostsService) GetPageSettings(ctx context.Context, req *pb.GetPageSettingsRequest) (*pb.GetPageSettingsResponse, error) {
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	pageSettings, err := conn.Pages().GetByAccountId(dom.AccountId(req.AccountId))

	if err != nil {
		return nil, err
	}

	return &pb.GetPageSettingsResponse{
		VisibleForUnauthorized: pageSettings.VisibleForUnauthorized,
		CommentsEnabled:        pageSettings.CommentsEnabled,
		AnyoneCanPost:          pageSettings.AnyoneCanPost,
	}, nil
}

func (s *PostsService) NewPost(ctx context.Context, req *pb.NewPostRequest) (*pb.NewPostResponse, error) {
	authorIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorIdVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	authorId := authorIdVal.(dom.AccountId)

	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	pageData, err := conn.Pages().GetByAccountId(dom.AccountId(req.PageAccountId))

	if err != nil {
		return nil, err
	}

	if authorId != dom.AccountId(req.PageAccountId) {
		if !pageData.AnyoneCanPost {
			return nil, status.Error(codes.PermissionDenied, "page owner prohibited posting")
		}
	}

	postId, err := conn.Posts().New(pageData.Id, repos.NewPostData{
		AuthorId: authorId,
		Content: dom.PostContent{
			Text:         dom.Text(req.Text),
			SourcePostId: opt.FromPointer((*dom.PostId)(req.Repost)),
		},
	})

	if err != nil {
		return nil, err
	}

	return &pb.NewPostResponse{
		PostId: int32(postId),
	}, nil
}

func (s *PostsService) NewComment(ctx context.Context, req *pb.NewCommentRequest) (*pb.NewCommentResponse, error) {
	authorIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorIdVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	authorId := authorIdVal.(dom.AccountId)

	commentsEnabledErr := func() error {
		conn, err := s.Db.OpenConnection(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()

		pageData, err := conn.Pages().GetByPostId(dom.PostId(req.PostId))
		if err != nil {
			return err
		}

		if !pageData.CommentsEnabled {
			return status.Error(codes.PermissionDenied, "comments prohibited")
		}

		return nil
	}()

	if commentsEnabledErr != nil {
		return nil, commentsEnabledErr
	}

	tx, err := s.Db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	commentId, err := tx.Comments().New(dom.PostId(req.PostId), repos.NewCommentData{
		AuthorId:       authorId,
		Content:        dom.Text(req.Content),
		ReplyCommentId: opt.FromPointer((*dom.CommentId)(req.ReplyCommentId)),
	})

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	payload, err := json.Marshal(statsModels.PostCommentEvent{
		CommentId:       statsModels.CommentId(commentId),
		PostId:          statsModels.PostId(req.PostId),
		AuthorAccountId: statsModels.AccountId(authorId),
		Timestamp:       time.Now(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Outbox().Put(dom.OutboxEvent{
		Type:    "comment",
		Payload: dom.OutboxEventPayload(payload),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &pb.NewCommentResponse{
		CommentId: int32(commentId),
	}, nil
}

func (s *PostsService) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil {
		pageData, err := conn.Pages().GetByPostId(dom.PostId(req.PostId))

		if err != nil {
			return nil, err
		}

		if !pageData.VisibleForUnauthorized {
			return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
		}
	}

	commentsList, err := conn.Comments().List(dom.PostId(req.PostId), repos.PagiToken(req.PageToken))
	if err != nil {
		return nil, err
	}

	comments := make([]*pb.Comment, len(commentsList.Comments))
	for i, comment := range commentsList.Comments {
		comments[i] = &pb.Comment{
			Id:              int32(comment.Id),
			AuthorAccountId: int32(comment.AuthorId),
			Content:         string(comment.Content),
			ReplyCommentId:  (*int32)(comment.ReplyId.ToPointer()),
		}
	}

	return &pb.GetCommentsResponse{
		Comments:      comments,
		NextPageToken: string(commentsList.NextPagiToken),
	}, nil
}

func (s *PostsService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	pageData, err := conn.Pages().GetByPostId(dom.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	post, err := conn.Posts().Get(dom.PostId(req.PostId))

	if err != nil {
		return nil, err
	}

	return &pb.Post{
		Id:              int32(post.Id),
		AuthorAccountId: int32(post.AuthorAccountId),
		Text:            string(post.Content.Text),
		SourcePostId:    (*int32)(post.Content.SourcePostId.ToPointer()),
		Pinned:          post.Pinned,
		ViewsCount:      post.ViewsCount,
	}, nil
}

func (s *PostsService) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	pageData, err := conn.Pages().GetByAccountId(dom.AccountId(req.PageAccountId))
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	postsList, err := conn.Posts().List(pageData.Id, repos.PagiToken(req.PageToken))
	if err != nil {
		return nil, err
	}

	posts := make([]*pb.Post, len(postsList.Posts))
	for i, post := range postsList.Posts {
		posts[i] = &pb.Post{
			Id:              int32(post.Id),
			AuthorAccountId: int32(post.AuthorAccountId),
			Text:            string(post.Content.Text),
			SourcePostId:    (*int32)(post.Content.SourcePostId.ToPointer()),
			Pinned:          post.Pinned,
			ViewsCount:      post.ViewsCount,
		}
	}

	return &pb.GetPostsResponse{
		Posts:         posts,
		NextPageToken: string(postsList.NextPagiToken),
	}, nil
}

func (s *PostsService) EditPost(ctx context.Context, req *pb.EditPostRequest) (*pb.Empty, error) {
	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	post, err := conn.Posts().Get(dom.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	if post.AuthorAccountId != authorizedId.(dom.AccountId) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if (req.Text == nil || *req.Text == string(post.Content.Text)) && (req.Pinned == nil || post.Pinned == *req.Pinned) {
		return &pb.Empty{}, nil
	}

	err = conn.Posts().Edit(post.Id, repos.EditedPostData{
		Text:   opt.FromPointer((*dom.Text)(req.Text)),
		Pinned: opt.FromPointer(req.Pinned),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *PostsService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.Empty, error) {
	authorizedIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedIdVal == nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	authorizedId := authorizedIdVal.(dom.AccountId)

	conn, err := s.Db.OpenConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	post, err := conn.Posts().Get(dom.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	if post.AuthorAccountId != authorizedId {
		page, err := conn.Pages().GetByPageId(dom.PageId(post.PageId))

		if err != nil {
			return nil, err
		}

		if page.AccountId != authorizedId {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}
	}

	err = conn.Posts().Delete(dom.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, err
}

func (s *PostsService) NewView(ctx context.Context, req *pb.NewViewRequest) (*pb.Empty, error) {
	authorizedIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedIdVal == nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	authorizedId := authorizedIdVal.(dom.AccountId)

	tx, err := s.Db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	err = tx.Metrics().NewView(authorizedId, dom.PostId(req.PostId))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	payload, err := json.Marshal(statsModels.PostViewEvent{
		PostId:          statsModels.PostId(req.PostId),
		ViewerAccountId: statsModels.AccountId(authorizedId),
		Timestamp:       time.Now(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Outbox().Put(dom.OutboxEvent{
		Type:    "view",
		Payload: dom.OutboxEventPayload(payload),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *PostsService) NewLike(ctx context.Context, req *pb.NewLikeRequest) (*pb.Empty, error) {
	authorizedIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedIdVal == nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	authorizedId := authorizedIdVal.(dom.AccountId)

	tx, err := s.Db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	err = tx.Metrics().NewLike(authorizedId, dom.PostId(req.PostId))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	payload, err := json.Marshal(statsModels.PostLikeEvent{
		PostId:         statsModels.PostId(req.PostId),
		LikerAccountId: statsModels.AccountId(authorizedId),
		Timestamp:      time.Now(),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Outbox().Put(dom.OutboxEvent{
		Type:    "like",
		Payload: dom.OutboxEventPayload(payload),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

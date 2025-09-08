package server

import (
	"context"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	opt "soa-socialnetwork/services/common/option"
	dom "soa-socialnetwork/services/posts/internal/domain"
	"soa-socialnetwork/services/posts/internal/repos"
	"soa-socialnetwork/services/posts/internal/server/interceptors"
	"soa-socialnetwork/services/posts/internal/storage/postgres"
	pb "soa-socialnetwork/services/posts/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repositories struct {
	Pages    repos.PagesRepository
	Posts    repos.PostsRepository
	Comments repos.CommentsRepository
	Metrics  repos.MetricsRepository
}

type PostsService struct {
	pb.UnimplementedPostsServiceServer

	storage     repositories
	jwtVerifier soajwt.Verifier
}

func newPostsService(cfg PostsServiceConfig) (PostsService, error) {
	pgPool, err := postgres.NewPool(postgres.ConnectionConfig{
		Host:     cfg.DbHost,
		User:     cfg.DbUser,
		Password: cfg.DbPassword,
		PoolSize: cfg.DbPoolSize,
	})

	if err != nil {
		return PostsService{}, err
	}

	return PostsService{
		storage: repositories{
			Pages:    &postgres.PageRepo{ConnPool: pgPool},
			Posts:    &postgres.PostsRepo{ConnPool: pgPool},
			Comments: &postgres.CommentsRepo{ConnPool: pgPool},
			Metrics:  &postgres.MetricsRepo{ConnPool: pgPool},
		},
		jwtVerifier: soajwt.NewVerifier(cfg.JwtPublicKey),
	}, nil
}

func (s *PostsService) EditPageSettings(ctx context.Context, req *pb.EditPageSettingsRequest) (*pb.Empty, error) {
	accountId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if accountId == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	if dom.AccountId(req.AccountId) != accountId.(dom.AccountId) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	pageData, err := s.storage.Pages.Get(ctx, repos.AccountId(req.AccountId))

	if err != nil {
		return nil, err
	}

	err = s.storage.Pages.Edit(ctx, pageData.Id, repos.EditedPageSettings{
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
	pageSettings, err := s.storage.Pages.Get(ctx, repos.AccountId(req.AccountId))

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

	pageData, err := s.storage.Pages.Get(ctx, repos.AccountId(req.PageAccountId))

	if err != nil {
		return nil, err
	}

	if authorId != dom.AccountId(req.PageAccountId) {
		if !pageData.AnyoneCanPost {
			return nil, status.Error(codes.PermissionDenied, "page owner prohibited posting")
		}
	}

	postId, err := s.storage.Posts.New(ctx, pageData.Id, repos.NewPostData{
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

	pageData, err := s.storage.Pages.Get(ctx, repos.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	if !pageData.CommentsEnabled {
		return nil, status.Error(codes.PermissionDenied, "comments prohibited")
	}

	commentId, err := s.storage.Comments.New(ctx, dom.PostId(req.PostId), repos.NewCommentData{
		AuthorId:       authorId,
		Content:        dom.Text(req.Content),
		ReplyCommentId: opt.FromPointer((*dom.CommentId)(req.ReplyCommentId)),
	})

	if err != nil {
		return nil, err
	}

	return &pb.NewCommentResponse{
		CommentId: int32(commentId),
	}, nil
}

func (s *PostsService) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil {
		pageData, err := s.storage.Pages.Get(ctx, repos.PostId(req.PostId))

		if err != nil {
			return nil, err
		}

		if !pageData.VisibleForUnauthorized {
			return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
		}
	}

	commentsList, err := s.storage.Comments.List(ctx, dom.PostId(req.PostId), repos.PagiToken(req.PageToken))
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
	pageData, err := s.storage.Pages.Get(ctx, repos.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	post, err := s.storage.Posts.Get(ctx, dom.PostId(req.PostId))

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
	pageData, err := s.storage.Pages.Get(ctx, repos.AccountId(req.PageAccountId))
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	postsList, err := s.storage.Posts.List(ctx, pageData.Id, repos.PagiToken(req.PageToken))
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

	post, err := s.storage.Posts.Get(ctx, dom.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	if post.AuthorAccountId != authorizedId.(dom.AccountId) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if (req.Text == nil || *req.Text == string(post.Content.Text)) && (req.Pinned == nil || post.Pinned == *req.Pinned) {
		return &pb.Empty{}, nil
	}

	err = s.storage.Posts.Edit(ctx, post.Id, repos.EditedPostData{
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

	post, err := s.storage.Posts.Get(ctx, dom.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	if post.AuthorAccountId != authorizedId {
		page, err := s.storage.Pages.Get(ctx, repos.PageId(post.PageId))

		if err != nil {
			return nil, err
		}

		if page.AccountId != authorizedId {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}
	}

	err = s.storage.Posts.Delete(ctx, dom.PostId(req.PostId))
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

	err := s.storage.Metrics.NewView(ctx, authorizedId, dom.PostId(req.PostId))
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

	err := s.storage.Metrics.NewLike(ctx, authorizedId, dom.PostId(req.PostId))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

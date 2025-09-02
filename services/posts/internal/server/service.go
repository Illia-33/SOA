package server

import (
	"context"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	db "soa-socialnetwork/services/posts/internal/server/dbclient"
	dbReq "soa-socialnetwork/services/posts/internal/server/dbclient/requests"
	dbt "soa-socialnetwork/services/posts/internal/server/dbclient/types"
	"soa-socialnetwork/services/posts/internal/server/interceptors"
	pb "soa-socialnetwork/services/posts/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsService struct {
	pb.UnimplementedPostsServiceServer

	dbClient    db.DatabaseClient
	jwtVerifier soajwt.Verifier
}

func newPostsService(cfg PostsServiceConfig) (PostsService, error) {
	dbc, err := db.NewPostgresDbClient(db.PostgresConfig{
		Host:     cfg.DbHost,
		User:     cfg.DbUser,
		Password: cfg.DbPassword,
		PoolSize: cfg.DbPoolSize,
	})

	if err != nil {
		return PostsService{}, err
	}

	return PostsService{
		dbClient:    dbc,
		jwtVerifier: soajwt.NewVerifier(cfg.JwtPublicKey),
	}, nil
}

func (s *PostsService) EditPageSettings(ctx context.Context, req *pb.EditPageSettingsRequest) (*pb.Empty, error) {
	accountId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if accountId == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	if dbt.AccountId(req.AccountId) != accountId.(dbt.AccountId) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	pageData, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
		EntityId: dbReq.AccountId(req.AccountId),
	})

	if err != nil {
		return nil, err
	}

	err = s.dbClient.EditPageSettings(ctx, dbReq.EditPageSettingsRequest{
		PageId:                 pageData.Id,
		VisibleForUnauthorized: dbt.OptionFromPtr(req.VisibleForUnauthorized),
		CommentsEnabled:        dbt.OptionFromPtr(req.CommentsEnabled),
		AnyoneCanPost:          dbt.OptionFromPtr(req.AnyoneCanPost),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *PostsService) GetPageSettings(ctx context.Context, req *pb.GetPageSettingsRequest) (*pb.GetPageSettingsResponse, error) {
	st, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
		EntityId: dbReq.AccountId(req.AccountId),
	})

	if err != nil {
		return nil, err
	}

	return &pb.GetPageSettingsResponse{
		VisibleForUnauthorized: st.VisibleForUnauthorized,
		CommentsEnabled:        st.CommentsEnabled,
		AnyoneCanPost:          st.AnyoneCanPost,
	}, nil
}

func (s *PostsService) NewPost(ctx context.Context, req *pb.NewPostRequest) (*pb.NewPostResponse, error) {
	authorIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorIdVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	authorId := authorIdVal.(dbt.AccountId)

	pageData, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
		EntityId: dbReq.AccountId(req.PageAccountId),
	})

	if err != nil {
		return nil, err
	}

	if authorId != dbt.AccountId(req.PageAccountId) {
		if !pageData.AnyoneCanPost {
			return nil, status.Error(codes.PermissionDenied, "page owner prohibited posting")
		}
	}

	dbResp, err := s.dbClient.NewPost(ctx, dbReq.NewPostRequest{
		PageId:   pageData.Id,
		AuthorId: authorId,
		Content: dbt.PostContent{
			Text:         dbt.Text(req.Text),
			SourcePostId: dbt.OptionFromPtr((*dbt.PostId)(req.Repost)),
		},
	})

	if err != nil {
		return nil, err
	}

	return &pb.NewPostResponse{
		PostId: int32(dbResp.Id),
	}, nil
}

func (s *PostsService) NewComment(ctx context.Context, req *pb.NewCommentRequest) (*pb.NewCommentResponse, error) {
	authorIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorIdVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	authorId := authorIdVal.(dbt.AccountId)

	pageData, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
		EntityId: dbReq.PostId(req.PostId),
	})
	if err != nil {
		return nil, err
	}
	if !pageData.CommentsEnabled {
		return nil, status.Error(codes.PermissionDenied, "comments prohibited")
	}

	dbResp, err := s.dbClient.NewComment(ctx, dbReq.NewCommentRequest{
		PostId:         dbt.PostId(req.PostId),
		AuthorId:       authorId,
		Content:        dbt.Text(req.Content),
		ReplyCommentId: dbt.OptionFromPtr((*dbt.CommentId)(req.ReplyCommentId)),
	})

	if err != nil {
		return nil, err
	}

	return &pb.NewCommentResponse{
		CommentId: int32(dbResp.Id),
	}, nil
}

func (s *PostsService) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil {
		pageData, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
			EntityId: dbReq.PostId(req.PostId),
		})

		if err != nil {
			return nil, err
		}

		if !pageData.VisibleForUnauthorized {
			return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
		}
	}

	dbResponse, err := s.dbClient.GetComments(ctx, dbReq.GetCommentsRequest{
		PostId:    dbt.PostId(req.PostId),
		PageToken: req.PageToken,
	})
	if err != nil {
		return nil, err
	}

	comments := make([]*pb.Comment, len(dbResponse.Comments))
	for i, comment := range dbResponse.Comments {
		comments[i] = &pb.Comment{
			Id:              int32(comment.Id),
			AuthorAccountId: int32(comment.AuthorId),
			Content:         string(comment.Content),
			ReplyCommentId:  (*int32)(comment.ReplyId.ToPointer()),
		}
	}

	return &pb.GetCommentsResponse{
		Comments:      comments,
		NextPageToken: string(dbResponse.NextPageToken),
	}, nil
}

func (s *PostsService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	pageData, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
		EntityId: dbReq.PostId(req.PostId),
	})
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	dbResp, err := s.dbClient.GetPost(ctx, dbReq.GetPostRequest{
		PostId: dbt.PostId(req.PostId),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Post{
		Id:              int32(dbResp.Post.Id),
		AuthorAccountId: int32(dbResp.Post.AuthorAccountId),
		Text:            string(dbResp.Post.Content.Text),
		SourcePostId:    (*int32)(dbResp.Post.Content.SourcePostId.ToPointer()),
		Pinned:          dbResp.Post.Pinned,
		ViewsCount:      dbResp.Post.ViewsCount,
	}, nil
}

func (s *PostsService) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	pageData, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
		EntityId: dbReq.AccountId(req.PageAccountId),
	})
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	dbResp, err := s.dbClient.GetPosts(ctx, dbReq.GetPostsRequest{
		PageId:    pageData.Id,
		PageToken: dbReq.PagiToken(req.PageToken),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]*pb.Post, len(dbResp.Posts))
	for i, post := range dbResp.Posts {
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
		NextPageToken: string(dbResp.NextPageToken),
	}, nil
}

func (s *PostsService) EditPost(ctx context.Context, req *pb.EditPostRequest) (*pb.Empty, error) {
	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	postData, err := s.dbClient.GetPost(ctx, dbReq.GetPostRequest{
		PostId: dbt.PostId(req.PostId),
	})
	if err != nil {
		return nil, err
	}

	if postData.Post.AuthorAccountId != authorizedId.(dbt.AccountId) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if (req.Text == nil || *req.Text == string(postData.Post.Content.Text)) && (req.Pinned == nil || postData.Post.Pinned == *req.Pinned) {
		return &pb.Empty{}, nil
	}

	err = s.dbClient.EditPost(ctx, dbReq.EditPostRequest{
		PostId: dbt.PostId(req.PostId),
		Text:   dbt.OptionFromPtr((*dbt.Text)(req.Text)),
		Pinned: dbt.OptionFromPtr(req.Pinned),
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
	authorizedId := authorizedIdVal.(dbt.AccountId)

	postData, err := s.dbClient.GetPost(ctx, dbReq.GetPostRequest{
		PostId: dbt.PostId(req.PostId),
	})
	if err != nil {
		return nil, err
	}

	if postData.Post.AuthorAccountId != authorizedId {
		page, err := s.dbClient.GetPageData(ctx, dbReq.GetPageDataRequest{
			EntityId: dbReq.PageId(postData.Post.PageId),
		})

		if err != nil {
			return nil, err
		}

		if page.AccountId != authorizedId {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}
	}

	err = s.dbClient.DeletePost(ctx, dbReq.DeletePostRequest{
		PostId: dbt.PostId(req.PostId),
	})
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
	authorizedId := authorizedIdVal.(dbt.AccountId)

	err := s.dbClient.NewView(ctx, dbReq.NewViewRequest{
		AccountId: authorizedId,
		PostId:    dbt.PostId(req.PostId),
	})

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
	authorizedId := authorizedIdVal.(dbt.AccountId)

	err := s.dbClient.NewLike(ctx, dbReq.NewLikeRequest{
		AccountId: authorizedId,
		PostId:    dbt.PostId(req.PostId),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

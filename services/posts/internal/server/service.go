package server

import (
	"context"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	db "soa-socialnetwork/services/posts/internal/server/dbclient"
	dbreq "soa-socialnetwork/services/posts/internal/server/dbclient/requests"
	dbt "soa-socialnetwork/services/posts/internal/server/dbclient/types"
	"soa-socialnetwork/services/posts/internal/server/interceptors"
	pb "soa-socialnetwork/services/posts/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsService struct {
	pb.UnimplementedPostsServiceServer

	dbCliennt   db.DatabaseClient
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
		dbCliennt:   dbc,
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

	pageData, err := s.dbCliennt.GetPageData(ctx, dbreq.GetPageDataRequest{
		EntityId: dbreq.AccountId(req.AccountId),
	})

	if err != nil {
		return nil, err
	}

	err = s.dbCliennt.EditPageSettings(ctx, dbreq.EditPageSettingsRequest{
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
	st, err := s.dbCliennt.GetPageData(ctx, dbreq.GetPageDataRequest{
		EntityId: dbreq.AccountId(req.AccountId),
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

	pageData, err := s.dbCliennt.GetPageData(ctx, dbreq.GetPageDataRequest{
		EntityId: dbreq.AccountId(req.PageAccountId),
	})

	if err != nil {
		return nil, err
	}

	if authorId != dbt.AccountId(req.PageAccountId) {
		if !pageData.AnyoneCanPost {
			return nil, status.Error(codes.PermissionDenied, "page owner prohibited posting")
		}
	}

	r, err := s.dbCliennt.NewPost(ctx, dbreq.NewPostRequest{
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
		PostId: int32(r.Id),
	}, nil
}

func (s *PostsService) NewComment(ctx context.Context, req *pb.NewCommentRequest) (*pb.NewCommentResponse, error) {
	authorIdVal := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorIdVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	authorId := authorIdVal.(dbt.AccountId)

	pageData, err := s.dbCliennt.GetPageData(ctx, dbreq.GetPageDataRequest{
		EntityId: dbreq.PostId(req.PostId),
	})
	if err != nil {
		return nil, err
	}
	if !pageData.CommentsEnabled {
		return nil, status.Error(codes.PermissionDenied, "comments prohibited")
	}

	r, err := s.dbCliennt.NewComment(ctx, dbreq.NewCommentRequest{
		PostId:         dbt.PostId(req.PostId),
		AuthorId:       authorId,
		Content:        dbt.Text(req.Content),
		ReplyCommentId: dbt.OptionFromPtr((*dbt.CommentId)(req.ReplyCommentId)),
	})

	if err != nil {
		return nil, err
	}

	return &pb.NewCommentResponse{
		CommentId: int32(r.Id),
	}, nil
}

func (s *PostsService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	pageData, err := s.dbCliennt.GetPageData(ctx, dbreq.GetPageDataRequest{
		EntityId: dbreq.PostId(req.PostId),
	})
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	r, err := s.dbCliennt.GetPost(ctx, dbreq.GetPostRequest{
		PostId: dbt.PostId(req.PostId),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Post{
		Id:              int32(r.Post.Id),
		AuthorAccountId: int32(r.Post.AuthorAccountId),
		Text:            string(r.Post.Content.Text),
		SourcePostId:    (*int32)(r.Post.Content.SourcePostId.ToPointer()),
		Pinned:          r.Post.Pinned,
	}, nil
}

func (s *PostsService) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	pageData, err := s.dbCliennt.GetPageData(ctx, dbreq.GetPageDataRequest{
		EntityId: dbreq.AccountId(req.PageAccountId),
	})
	if err != nil {
		return nil, err
	}

	authorizedId := ctx.Value(interceptors.AUTHOR_ACCOUNT_ID_CTX_KEY)
	if authorizedId == nil && !pageData.VisibleForUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "denied for unauthorized")
	}

	r, err := s.dbCliennt.GetPosts(ctx, dbreq.GetPostsRequest{
		PageId:    pageData.Id,
		PageToken: dbreq.PaginationToken(req.PageToken),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]*pb.Post, len(r.Posts))
	for i, post := range r.Posts {
		posts[i] = &pb.Post{
			Id:              int32(post.Id),
			AuthorAccountId: int32(post.AuthorAccountId),
			Text:            string(post.Content.Text),
			SourcePostId:    (*int32)(post.Content.SourcePostId.ToPointer()),
			Pinned:          post.Pinned,
		}
	}

	return &pb.GetPostsResponse{
		Posts:         posts,
		NextPageToken: string(r.NextPageToken),
	}, nil
}

package server

import (
	"context"
	"fmt"
	"net/http"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	accountsPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/httperr"
	"soa-socialnetwork/services/gateway/internal/server/query"
	"soa-socialnetwork/services/gateway/internal/server/soagrpc"
	"soa-socialnetwork/services/gateway/pkg/types"
	postsPb "soa-socialnetwork/services/posts/proto"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GatewayService struct {
	jwtVerifier         soajwt.Verifier
	accountsGrpcTarget  string
	accountsStubFactory soagrpc.AccountsStubFactory
	postsStubFactory    soagrpc.PostsStubFactory
}

func newGatewayService(cfg GatewayServiceConfig) GatewayService {
	if cfg.AccountsServiceStubFactory == nil {
		cfg.AccountsServiceStubFactory = defaultAccountsStubFactory{}
	}

	if cfg.PostsServiceStubFactory == nil {
		cfg.PostsServiceStubFactory = defaultPostsStubFactory{}
	}

	return GatewayService{
		jwtVerifier:         soajwt.NewVerifier(cfg.JwtPublicKey),
		accountsGrpcTarget:  fmt.Sprintf("%s:%d", cfg.AccountsServiceHost, cfg.AccountsServicePort),
		accountsStubFactory: cfg.AccountsServiceStubFactory,
		postsStubFactory:    cfg.PostsServiceStubFactory,
	}
}

func (s *GatewayService) createAccountsStub(qp *query.Params) (accountsPb.AccountsServiceClient, error) {
	return s.accountsStubFactory.New(s.accountsGrpcTarget, qp)
}

func (s *GatewayService) createPostsStub(qp *query.Params) (postsPb.PostsServiceClient, error) {
	return s.postsStubFactory.New(s.accountsGrpcTarget, qp)
}

func (s *GatewayService) RegisterProfile(qp *query.Params, req *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.RegisterProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.RegisterUser(context.Background(), &accountsPb.RegisterUserRequest{
		Login:       string(req.Login),
		Password:    string(req.Password),
		Email:       string(req.Email),
		PhoneNumber: string(req.PhoneNumber),
		Name:        string(req.Name),
		Surname:     string(req.Surname),
	})
	if err != nil {
		return api.RegisterProfileResponse{}, httperr.FromGrpcError(err)
	}

	return api.RegisterProfileResponse{
		ProfileId: resp.ProfileId,
	}, httperr.Ok()
}

func (s *GatewayService) GetProfileInfo(qp *query.Params) (api.GetProfileResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.GetProfileResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetProfile(context.Background(), &accountsPb.GetProfileRequest{
		ProfileId: qp.ProfileId,
	})
	if err != nil {
		return api.GetProfileResponse{}, httperr.FromGrpcError(err)
	}

	return api.GetProfileResponse{
		Name:     resp.Name,
		Surname:  resp.Surname,
		Birthday: resp.Birthday.AsTime().Format("2006-01-02"),
		Bio:      resp.Bio,
	}, httperr.Ok()
}

func (s *GatewayService) EditProfileInfo(qp *query.Params, req *api.EditProfileRequest) httperr.Err {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	var pbBirthday *timestamppb.Timestamp = nil
	if req.Birthday.HasValue {
		pbBirthday = timestamppb.New(req.Birthday.Value.Time)
	}

	_, err = stub.EditProfile(context.Background(), &accountsPb.EditProfileRequest{
		ProfileId: qp.ProfileId,
		EditedProfileData: &accountsPb.Profile{
			Name:     string(req.Name.Value),
			Surname:  string(req.Surname.Value),
			Birthday: pbBirthday,
			Bio:      string(req.Bio.Value),
		},
	})
	if err != nil {
		return httperr.FromGrpcError(err)
	}

	return httperr.Ok()
}

func (s *GatewayService) DeleteProfile(qp *query.Params) httperr.Err {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.UnregisterUser(context.Background(), &accountsPb.UnregisterUserRequest{
		ProfileId: qp.ProfileId,
	})
	if err != nil {
		return httperr.FromGrpcError(err)
	}

	return httperr.Ok()
}

func (s *GatewayService) buildAuthByPassword(req *api.AuthenticateRequest) (proto accountsPb.AuthByPassword) {
	if req.Login.HasValue {
		proto.UserId = &accountsPb.AuthByPassword_Login{
			Login: string(req.Login.Value),
		}
	} else if req.Email.HasValue {
		proto.UserId = &accountsPb.AuthByPassword_Email{
			Email: string(req.Email.Value),
		}
	} else if req.PhoneNumber.HasValue {
		proto.UserId = &accountsPb.AuthByPassword_PhoneNumber{
			PhoneNumber: string(req.PhoneNumber.Value),
		}
	} else {
		panic("at least one user id must be provided")
	}

	proto.Password = string(req.Password)
	return
}

func (s *GatewayService) Authenticate(qp *query.Params, req *api.AuthenticateRequest) (api.AuthenticateResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.AuthenticateResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	protoRequest := s.buildAuthByPassword(req)
	resp, err := stub.Authenticate(context.Background(), &protoRequest)

	if err != nil {
		return api.AuthenticateResponse{}, httperr.FromGrpcError(err)
	}

	return api.AuthenticateResponse{
		Token: resp.Token,
	}, httperr.Ok()
}

func (s *GatewayService) CreateApiToken(qp *query.Params, req *api.CreateApiTokenRequest) (api.CreateApiTokenResponse, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.CreateApiTokenResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	protoAuthByPassword := s.buildAuthByPassword(&req.Auth)

	resp, err := stub.CreateApiToken(context.Background(), &accountsPb.CreateApiTokenRequest{
		Auth: &protoAuthByPassword,
		Params: &accountsPb.AuthTokenParams{
			ReadAccess:  req.ReadAccess,
			WriteAccess: req.WriteAccess,
			Ttl:         durationpb.New(req.Ttl.Duration),
		},
	})

	if err != nil {
		return api.CreateApiTokenResponse{}, httperr.FromGrpcError(err)
	}

	return api.CreateApiTokenResponse{
		Token: resp.Token,
	}, httperr.Ok()
}

func (s *GatewayService) resolveProfileId(qp *query.Params, profileId string) (int32, httperr.Err) {
	stub, err := s.createAccountsStub(qp)
	if err != nil {
		return 0, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.ResolveProfileId(context.Background(), &accountsPb.ResolveProfileIdRequest{
		ProfileId: profileId,
	})
	if err != nil {
		return 0, httperr.FromGrpcError(err)
	}

	return resp.AccountId, httperr.Ok()
}

func (s *GatewayService) GetPageSettings(qp *query.Params) (api.GetPageSettingsResponse, httperr.Err) {
	accountId, accErr := s.resolveProfileId(qp, qp.ProfileId)
	if !accErr.IsOk() {
		return api.GetPageSettingsResponse{}, accErr
	}

	stub, err := s.createPostsStub(qp)
	if err != nil {
		return api.GetPageSettingsResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetPageSettings(context.Background(), &postsPb.GetPageSettingsRequest{
		AccountId: accountId,
	})
	if err != nil {
		return api.GetPageSettingsResponse{}, httperr.FromGrpcError(err)
	}

	return api.GetPageSettingsResponse{
		VisibleForUnauthorized: resp.VisibleForUnauthorized,
		CommentsEnabled:        resp.CommentsEnabled,
		AnyoneCanPost:          resp.AnyoneCanPost,
	}, httperr.Ok()
}

func (s *GatewayService) EditPageSettings(qp *query.Params, req *api.EditPageSettingsRequest) httperr.Err {
	accountId, accErr := s.resolveProfileId(qp, qp.ProfileId)
	if !accErr.IsOk() {
		return accErr
	}

	stub, err := s.createPostsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.EditPageSettings(context.Background(), &postsPb.EditPageSettingsRequest{
		AccountId:              accountId,
		VisibleForUnauthorized: req.VisibleForUnauthorized.ToPointer(),
		CommentsEnabled:        req.CommentsEnabled.ToPointer(),
		AnyoneCanPost:          req.AnyoneCanPost.ToPointer(),
	})
	if err != nil {
		return httperr.FromGrpcError(err)
	}

	return httperr.Ok()
}

func (s *GatewayService) NewPost(qp *query.Params, req *api.NewPostRequest) (api.NewPostResponse, httperr.Err) {
	accountId, accErr := s.resolveProfileId(qp, qp.ProfileId)
	if !accErr.IsOk() {
		return api.NewPostResponse{}, accErr
	}

	stub, err := s.createPostsStub(qp)
	if err != nil {
		return api.NewPostResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.NewPost(context.Background(), &postsPb.NewPostRequest{
		PageAccountId: accountId,
		Text:          req.Text,
		Repost:        req.Repost.ToPointer(),
	})
	if err != nil {
		return api.NewPostResponse{}, httperr.FromGrpcError(err)
	}

	return api.NewPostResponse{
		PostId: resp.PostId,
	}, httperr.Ok()
}

func postFromProto(proto *postsPb.Post) api.Post {
	return api.Post{
		Id:           proto.Id,
		AuthorId:     proto.AuthorAccountId,
		Text:         proto.Text,
		SourcePostId: types.OptionalFromPointer(proto.SourcePostId),
		Pinned:       proto.Pinned,
	}
}

func (s *GatewayService) GetPost(qp *query.Params) (api.Post, httperr.Err) {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return api.Post{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetPost(context.Background(), &postsPb.GetPostRequest{
		PostId: qp.PostId,
	})
	if err != nil {
		return api.Post{}, httperr.FromGrpcError(err)
	}

	return postFromProto(resp), httperr.Ok()
}

func (s *GatewayService) GetPosts(qp *query.Params, req *api.GetPostsRequest) (api.GetPostsResponse, httperr.Err) {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return api.GetPostsResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	accountId, accErr := s.resolveProfileId(qp, qp.ProfileId)
	if !accErr.IsOk() {
		return api.GetPostsResponse{}, accErr
	}

	resp, err := stub.GetPosts(context.Background(), &postsPb.GetPostsRequest{
		PageAccountId: accountId,
	})
	if err != nil {
		return api.GetPostsResponse{}, httperr.FromGrpcError(err)
	}

	posts := make([]api.Post, len(resp.Posts))
	for i, p := range resp.Posts {
		posts[i] = postFromProto(p)
	}

	return api.GetPostsResponse{
		Posts:         posts,
		NextPageToken: resp.NextPageToken,
	}, httperr.Ok()
}

func (s *GatewayService) EditPost(qp *query.Params, req *api.EditPostRequest) httperr.Err {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.EditPost(context.Background(), &postsPb.EditPostRequest{
		PostId: qp.PostId,
		Text:   req.Text.ToPointer(),
		Pinned: req.Pinned.ToPointer(),
	})
	if err != nil {
		return httperr.FromGrpcError(err)
	}

	return httperr.Ok()
}

func (s *GatewayService) DeletePost(qp *query.Params) httperr.Err {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.DeletePost(context.Background(), &postsPb.DeletePostRequest{
		PostId: qp.PostId,
	})
	if err != nil {
		return httperr.FromGrpcError(err)
	}

	return httperr.Ok()
}

func (s *GatewayService) NewComment(qp *query.Params, req *api.NewCommentRequest) (api.NewCommentResponse, httperr.Err) {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return api.NewCommentResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.NewComment(context.Background(), &postsPb.NewCommentRequest{
		PostId:         qp.PostId,
		Content:        req.Content,
		ReplyCommentId: req.ReplyCommentId.ToPointer(),
	})
	if err != nil {
		return api.NewCommentResponse{}, httperr.FromGrpcError(err)
	}

	return api.NewCommentResponse{
		CommentId: resp.CommentId,
	}, httperr.Ok()
}

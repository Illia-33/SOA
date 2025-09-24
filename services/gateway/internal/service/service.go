package service

import (
	"context"
	"fmt"
	"net/http"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	accountsPb "soa-socialnetwork/services/accounts/proto"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/grpcutils"
	"soa-socialnetwork/services/gateway/internal/httperr"
	"soa-socialnetwork/services/gateway/internal/query"
	"soa-socialnetwork/services/gateway/pkg/types"
	postsPb "soa-socialnetwork/services/posts/proto"
	statsPb "soa-socialnetwork/services/stats/proto"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GatewayService struct {
	JwtVerifier          soajwt.Verifier
	AccountsGrpcAccessor GrpcAccessor[accountsPb.AccountsServiceClient]
	PostsGrpcAccessor    GrpcAccessor[postsPb.PostsServiceClient]
	StatsGrpcAccessor    GrpcAccessor[statsPb.StatsServiceClient]
}

type GrpcAccessor[TStub any] struct {
	Target  string
	Factory grpcutils.StubCreator[TStub]
}

func (a *GrpcAccessor[TStub]) createStub(qp *query.Params) (TStub, error) {
	return a.Factory.New(a.Target, qp)
}

func NewGatewayService(cfg Config) GatewayService {
	return GatewayService{
		JwtVerifier: soajwt.NewVerifier(cfg.JwtPublicKey),
		AccountsGrpcAccessor: GrpcAccessor[accountsPb.AccountsServiceClient]{
			Target:  fmt.Sprintf("%s:%d", cfg.AccountsServiceHost, cfg.AccountsServicePort),
			Factory: defaultAccountsStubCreator{},
		},
		PostsGrpcAccessor: GrpcAccessor[postsPb.PostsServiceClient]{
			Target:  fmt.Sprintf("%s:%d", cfg.PostsServiceHost, cfg.PostsServicePort),
			Factory: defaultPostsStubCreator{},
		},
		StatsGrpcAccessor: GrpcAccessor[statsPb.StatsServiceClient]{
			Target:  fmt.Sprintf("%s:%d", cfg.StatsServiceHost, cfg.StatsServicePort),
			Factory: defaultStatsStubCreator{},
		},
	}
}

func (s *GatewayService) createAccountsStub(qp *query.Params) (accountsPb.AccountsServiceClient, error) {
	return s.AccountsGrpcAccessor.createStub(qp)
}

func (s *GatewayService) createPostsStub(qp *query.Params) (postsPb.PostsServiceClient, error) {
	return s.PostsGrpcAccessor.createStub(qp)
}

func (s *GatewayService) createStatsStub(qp *query.Params) (statsPb.StatsServiceClient, error) {
	return s.StatsGrpcAccessor.createStub(qp)
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
		ViewsCount:   proto.ViewsCount,
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
		PageToken:     req.PageToken,
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

func commentFromProto(proto *postsPb.Comment) api.Comment {
	return api.Comment{
		Id:             proto.Id,
		AuthorId:       proto.AuthorAccountId,
		Content:        proto.Content,
		ReplyCommentId: types.OptionalFromPointer(proto.ReplyCommentId),
	}
}

func (s *GatewayService) GetComments(qp *query.Params, req *api.GetCommentsRequest) (api.GetCommentsResponse, httperr.Err) {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return api.GetCommentsResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetComments(context.Background(), &postsPb.GetCommentsRequest{
		PostId:    qp.PostId,
		PageToken: req.PageToken,
	})
	if err != nil {
		return api.GetCommentsResponse{}, httperr.FromGrpcError(err)
	}

	comments := make([]api.Comment, len(resp.Comments))
	for i, comment := range resp.Comments {
		comments[i] = commentFromProto(comment)
	}

	return api.GetCommentsResponse{
		Comments:      comments,
		NextPageToken: resp.NextPageToken,
	}, httperr.Ok()
}

func (s *GatewayService) NewView(qp *query.Params) httperr.Err {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.NewView(context.Background(), &postsPb.NewViewRequest{
		PostId: qp.PostId,
	})
	if err != nil {
		return httperr.FromGrpcError(err)
	}

	return httperr.Ok()
}

func (s *GatewayService) NewLike(qp *query.Params) httperr.Err {
	stub, err := s.createPostsStub(qp)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, err)
	}

	_, err = stub.NewLike(context.Background(), &postsPb.NewLikeRequest{
		PostId: qp.PostId,
	})
	if err != nil {
		return httperr.FromGrpcError(err)
	}

	return httperr.Ok()
}

func (s *GatewayService) GetPostMetric(qp *query.Params, req *api.GetPostMetricRequest) (api.GetPostMetricResponse, httperr.Err) {
	stub, err := s.createStatsStub(qp)
	if err != nil {
		return api.GetPostMetricResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetPostMetric(context.Background(), &statsPb.GetPostMetricRequest{
		PostId: qp.PostId,
		Metric: metricToProto(req.Metric),
	})
	if err != nil {
		return api.GetPostMetricResponse{}, httperr.FromGrpcError(err)
	}

	return api.GetPostMetricResponse{
		Count: int(resp.Count),
	}, httperr.Ok()
}

func (s *GatewayService) GetPostMetricDynamics(qp *query.Params, req *api.GetPostMetricDynamicsRequest) (api.GetPostMetricDynamicsResponse, httperr.Err) {
	stub, err := s.createStatsStub(qp)
	if err != nil {
		return api.GetPostMetricDynamicsResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetPostMetricDynamics(context.Background(), &statsPb.GetPostMetricDynamicsRequest{
		PostId: qp.PostId,
		Metric: metricToProto(req.Metric),
	})
	if err != nil {
		return api.GetPostMetricDynamicsResponse{}, httperr.FromGrpcError(err)
	}

	dynamics := make([]api.DayDynamics, len(resp.Dynamics))
	for i, dyn := range resp.Dynamics {
		dynamics[i] = dayStatsFromProto(dyn)
	}

	return api.GetPostMetricDynamicsResponse{
		Dynamics: dynamics,
	}, httperr.Ok()
}

func (s *GatewayService) GetTop10Posts(qp *query.Params, req *api.GetTop10PostsRequest) (api.GetTop10PostsResponse, httperr.Err) {
	stub, err := s.createStatsStub(qp)
	if err != nil {
		return api.GetTop10PostsResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := stub.GetTop10Posts(context.Background(), &statsPb.GetTop10PostsRequest{
		Metric: metricToProto(req.Metric),
	})
	if err != nil {
		return api.GetTop10PostsResponse{}, httperr.FromGrpcError(err)
	}

	posts := make([]api.PostStats, len(resp.Posts))
	for i, postPb := range resp.Posts {
		posts[i] = postStatsFromProto(postPb, req.Metric)
	}

	return api.GetTop10PostsResponse{
		Posts: posts,
	}, httperr.Ok()
}

func (s *GatewayService) GetTop10Users(qp *query.Params, req *api.GetTop10UsersRequest) (api.GetTop10UsersResponse, httperr.Err) {
	ctx := context.Background()

	statsStub, err := s.createStatsStub(qp)
	if err != nil {
		return api.GetTop10UsersResponse{}, httperr.New(http.StatusInternalServerError, err)
	}

	resp, err := statsStub.GetTop10Users(ctx, &statsPb.GetTop10UsersRequest{
		Metric: metricToProto(req.Metric),
	})
	if err != nil {
		return api.GetTop10UsersResponse{}, httperr.FromGrpcError(err)
	}

	accountsStub, err := s.createAccountsStub(qp)
	if err != nil {
		return api.GetTop10UsersResponse{}, httperr.Ok()
	}

	users := make([]api.UserStats, len(resp.Users))
	for i, userPb := range resp.Users {
		resp, err := accountsStub.ResolveAccountId(ctx, &accountsPb.ResolveAccountIdRequest{
			AccountId: userPb.UserId,
		})
		if err != nil {
			return api.GetTop10UsersResponse{}, httperr.FromGrpcError(err)
		}

		users[i] = userStatsFromProto(userPb, req.Metric)
		users[i].Id = resp.ProfileId
	}

	return api.GetTop10UsersResponse{
		Users: users,
	}, httperr.Ok()
}

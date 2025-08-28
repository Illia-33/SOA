package api

import "soa-socialnetwork/services/gateway/pkg/types"

type RegisterProfileRequest struct {
	Login       types.Login       `json:"login"`
	Password    types.Password    `json:"password"`
	Email       types.Email       `json:"email"`
	PhoneNumber types.PhoneNumber `json:"phone_number"`
	Name        types.Name        `json:"name"`
	Surname     types.Surname     `json:"surname"`
}

type RegisterProfileResponse struct {
	ProfileId string `json:"profile_id"`
}

type GetProfileResponse struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Birthday string `json:"birthday"`
	Bio      string `json:"bio"`
}

type EditProfileRequest struct {
	Name     types.Optional[types.Name]    `json:"name"`
	Surname  types.Optional[types.Surname] `json:"surname"`
	Birthday types.Optional[types.Date]    `json:"birthday"`
	Bio      types.Optional[types.Bio]     `json:"bio"`
}

// Exactly one of {login,email,phone_number} must have value
type AuthenticateRequest struct {
	Login       types.Optional[types.Login]       `json:"login"`
	Email       types.Optional[types.Email]       `json:"email"`
	PhoneNumber types.Optional[types.PhoneNumber] `json:"phone_number"`
	Password    types.Password                    `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}

// Ttl (time to live) must be positive
type CreateApiTokenRequest struct {
	Auth        AuthenticateRequest `json:"auth"`
	ReadAccess  bool                `json:"read_access"`
	WriteAccess bool                `json:"write_access"`
	Ttl         types.Duration      `json:"ttl"`
}

type CreateApiTokenResponse struct {
	Token string `json:"token"`
}

type GetPageSettingsResponse struct {
	VisibleForUnauthorized bool `json:"visible_for_unauthorized"`
	CommentsEnabled        bool `json:"comments_enabled"`
	AnyoneCanPost          bool `json:"anyone_can_post"`
}

type EditPageSettingsRequest struct {
	VisibleForUnauthorized types.Optional[bool] `json:"visible_for_unauthorized"`
	CommentsEnabled        types.Optional[bool] `json:"comments_enabled"`
	AnyoneCanPost          types.Optional[bool] `json:"anyone_can_post"`
}

type NewPostRequest struct {
	Text   string                `json:"text"`
	Repost types.Optional[int32] `json:"repost"`
}

type NewPostResponse struct {
	PostId int32 `json:"post_id"`
}

type NewCommentRequest struct {
	Content        string                `json:"content"`
	ReplyCommentId types.Optional[int32] `json:"reply_comment_id"`
}

type NewCommentResponse struct {
	CommentId int32 `json:"comment_id"`
}

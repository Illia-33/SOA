package api

import "soa-socialnetwork/services/gateway/pkg/types"

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

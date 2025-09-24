package api

import "soa-socialnetwork/services/gateway/pkg/types"

type GetPostMetricRequest struct {
	Metric types.Metric `json:"metric"`
}

type GetPostMetricResponse struct {
	Count int `json:"count"`
}

type GetPostMetricDynamicsRequest struct {
	Metric types.Metric `json:"metric"`
}

type DayDynamics struct {
	Date  types.Date `json:"day"`
	Count int        `json:"count"`
}

type GetPostMetricDynamicsResponse struct {
	Dynamics []DayDynamics `json:"dynamics"`
}

type GetTop10PostsRequest struct {
	Metric types.Metric `json:"metric"`
}

type PostStats struct {
	Id    int `json:"post_id"`
	Value int `json:"value"`
}

type GetTop10PostsResponse struct {
	Posts []PostStats `json:"posts"`
}

type GetTop10UsersRequest struct {
	Metric types.Metric `json:"metric"`
}

type UserStats struct {
	Id    string `json:"user_id"`
	Value int    `json:"value"`
}

type GetTop10UsersResponse struct {
	Users []UserStats `json:"users"`
}

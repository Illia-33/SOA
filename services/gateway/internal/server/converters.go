package server

import (
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/pkg/types"
	statsPb "soa-socialnetwork/services/stats/proto"
)

func metricToProto(metric types.Metric) statsPb.Metric {
	switch metric {
	case types.METRIC_VIEW_COUNT:
		return statsPb.Metric_METRIC_VIEW_COUNT

	case types.METRIC_LIKE_COUNT:
		return statsPb.Metric_METRIC_LIKE_COUNT

	case types.METRIC_COMMENT_COUNT:
		return statsPb.Metric_METRIC_COMMENT_COUNT
	}

	return statsPb.Metric_METRIC_UNSPECIFIED
}

func dayStatsFromProto(dayStats *statsPb.DayDynamics) api.DayDynamics {
	return api.DayDynamics{
		Date: types.Date{
			Time: dayStats.Date.AsTime(),
		},
		Count: int(dayStats.Count),
	}
}

func postStatsFromProto(postStats *statsPb.PostStats, metric types.Metric) api.PostStats {
	var result api.PostStats
	result.Id = int(postStats.PostId)

	switch metric {
	case types.METRIC_VIEW_COUNT:
		result.Value = int(postStats.ViewCount)

	case types.METRIC_LIKE_COUNT:
		result.Value = int(postStats.LikeCount)

	case types.METRIC_COMMENT_COUNT:
		result.Value = int(postStats.CommentCount)
	}

	return result
}

func userStatsFromProto(userStats *statsPb.UserStats, metric types.Metric) api.UserStats {
	var result api.UserStats
	switch metric {
	case types.METRIC_VIEW_COUNT:
		result.Value = int(userStats.ViewCount)

	case types.METRIC_LIKE_COUNT:
		result.Value = int(userStats.LikeCount)

	case types.METRIC_COMMENT_COUNT:
		result.Value = int(userStats.CommentCount)
	}

	return result
}

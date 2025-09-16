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

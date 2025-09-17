package service

import (
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"
	pb "soa-socialnetwork/services/stats/proto"
)

func metricFromProto(value pb.Metric) models.Metric {
	switch value {
	case pb.Metric_METRIC_VIEW_COUNT:
		return models.METRIC_VIEW_COUNT

	case pb.Metric_METRIC_LIKE_COUNT:
		return models.METRIC_LIKE_COUNT

	case pb.Metric_METRIC_COMMENT_COUNT:
		return models.METRIC_COMMENT_COUNT
	}

	panic("unknown proto metric")
}

func postStatsFromAgg(value repo.PostAgg, metric models.Metric) *pb.PostStats {
	postStats := &pb.PostStats{
		PostId: int32(value.PostId),
	}

	switch metric {
	case models.METRIC_VIEW_COUNT:
		{
			postStats.ViewCount = uint64(value.MetricValue)
		}

	case models.METRIC_LIKE_COUNT:
		{
			postStats.LikeCount = uint64(value.MetricValue)
		}

	case models.METRIC_COMMENT_COUNT:
		{
			postStats.CommentCount = uint64(value.MetricValue)
		}

	default:
		panic("unknown metric")
	}

	return postStats
}

func userStatsFromAgg(value repo.UserAgg, metric models.Metric) *pb.UserStats {
	postStats := &pb.UserStats{
		UserId: int32(value.AccountId),
	}

	switch metric {
	case models.METRIC_VIEW_COUNT:
		{
			postStats.ViewCount = uint64(value.MetricValue)
		}

	case models.METRIC_LIKE_COUNT:
		{
			postStats.LikeCount = uint64(value.MetricValue)
		}

	case models.METRIC_COMMENT_COUNT:
		{
			postStats.CommentCount = uint64(value.MetricValue)
		}

	default:
		panic("unknown metric")
	}

	return postStats
}

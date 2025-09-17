package repo

import "soa-socialnetwork/services/stats/pkg/models"

type UserAgg struct {
	AccountId   int32
	MetricValue uint64
}

type PostAgg struct {
	PostId      int32
	MetricValue uint64
}

type AggregationRepo interface {
	GetTop10UsersByMetric(models.Metric) ([]UserAgg, error)
	GetTop10PostsByMetric(models.Metric) ([]PostAgg, error)
}

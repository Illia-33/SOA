package repo

import "soa-socialnetwork/services/stats/pkg/models"

type UserAgg struct {
	AccountId   models.AccountId
	MetricValue int
}

type PostAgg struct {
	PostId      models.PostId
	MetricValue int
}

type AggregationRepo interface {
	GetTop10UsersByMetric(models.Metric) ([]UserAgg, error)
	GetTop10PostsByMetric(models.Metric) ([]PostAgg, error)
}

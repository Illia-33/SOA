package clickhouse

import (
	"context"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"

	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type aggregationRepo struct {
	ctx  context.Context
	conn chDriver.Conn
}

func toClickhouseMetricName(metric models.Metric) string {
	switch metric {
	case models.METRIC_VIEW_COUNT:
		return "view_count"

	case models.METRIC_LIKE_COUNT:
		return "like_count"

	case models.METRIC_COMMENT_COUNT:
		return "comment_count"
	}

	return ""
}

func (r *aggregationRepo) GetTop10UsersByMetric(metric models.Metric) ([]repo.UserAgg, error) {
	sql := `
	SELECT account_id, cnt
	FROM agg_user_metrics
	WHERE metric = ?
	ORDER BY cnt DESC
	LIMIT 10;
	`

	rows, err := r.conn.Query(r.ctx, sql, toClickhouseMetricName(metric))
	if err != nil {
		return nil, err
	}

	var aggs []repo.UserAgg
	for rows.Next() {
		var userAgg repo.UserAgg

		err = rows.Scan(&userAgg.AccountId, &userAgg.MetricValue)
		if err != nil {
			return nil, err
		}

		aggs = append(aggs, userAgg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return aggs, nil
}
func (r *aggregationRepo) GetTop10PostsByMetric(metric models.Metric) ([]repo.PostAgg, error) {
	sql := `
	SELECT post_id, cnt
	FROM agg_post_metrics
	WHERE metric = ?
	ORDER BY cnt DESC
	LIMIT 10;
	`

	rows, err := r.conn.Query(r.ctx, sql, toClickhouseMetricName(metric))
	if err != nil {
		return nil, err
	}

	var aggs []repo.PostAgg
	for rows.Next() {
		var userAgg repo.PostAgg

		err = rows.Scan(&userAgg.PostId, &userAgg.MetricValue)
		if err != nil {
			return nil, err
		}

		aggs = append(aggs, userAgg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return aggs, nil
}

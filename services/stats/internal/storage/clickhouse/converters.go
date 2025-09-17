package clickhouse

import (
	"soa-socialnetwork/services/stats/pkg/models"
	"time"
)

func dateFromYYYYMMDD(ymd uint32) time.Time {
	ymdInt := int(ymd)
	return time.Date(ymdInt/10000, time.Month((ymdInt%10000)/100), (ymdInt % 1000000), 0, 0, 0, 0, time.UTC)
}

func metricToChMetricName(metric models.Metric) string {
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

package types

import (
	"encoding/json"
	"fmt"
)

type Metric int

const (
	METRIC_VIEW_COUNT Metric = iota
	METRIC_LIKE_COUNT
	METRIC_COMMENT_COUNT
)

func (n *Metric) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	metric, ok := string_to_metric[s]
	if !ok {
		return errorUnknownMetric()
	}

	*n = metric
	return nil
}

func (n Metric) MarshalJSON() ([]byte, error) {
	s, ok := metric_to_string[n]
	if !ok {
		return nil, fmt.Errorf("unknown metric: %v", n)
	}

	return json.Marshal(s)
}

var string_to_metric map[string]Metric = map[string]Metric{
	"view_count":    METRIC_VIEW_COUNT,
	"like_count":    METRIC_LIKE_COUNT,
	"comment_count": METRIC_COMMENT_COUNT,
}

var metric_to_string map[Metric]string = genMetricToString()

func genMetricToString() map[Metric]string {
	m := make(map[Metric]string)
	for k, v := range string_to_metric {
		m[v] = k
	}
	return m
}

func errorUnknownMetric() error {
	var possibleMetrics []string
	for key := range string_to_metric {
		possibleMetrics = append(possibleMetrics, key)
	}

	return fmt.Errorf("unknown metric, must be one of %v", possibleMetrics)
}

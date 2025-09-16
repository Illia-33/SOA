package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func tryGetMetric(t *testing.T, postId int, metric string) *http.Response {
	body := map[string]any{
		"metric": metric,
	}
	return makeRequest(t, http.MethodGet, fmt.Sprintf("/post/%d/metric", postId), body, "")
}

func getMetricOk(t *testing.T, postId int, metric string) int {
	resp := tryGetMetric(t, postId, metric)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return int(responseBodyToMap(t, resp)["count"].(float64))
}

func TestGetViewCount(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "get_view_count",
		"password":     "testpasswd",
		"email":        "get_view_count@yahoo.com",
		"phone_number": "+79250000022",
		"name":         "Get",
		"surname":      "ViewCount",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "get_view_count",
		"password": "testpasswd",
	})

	post := map[string]any{
		"text": "test post with views",
	}
	postId := createPostOk(t, id, post, jwtAuth(token))

	const VIEW_COUNT = 5
	for range VIEW_COUNT {
		newViewOk(t, postId, jwtAuth(token))
	}

	time.Sleep(5 * time.Second)

	viewCount := getMetricOk(t, postId, "view_count")
	require.Equal(t, VIEW_COUNT, viewCount)
}

func TestGetLikeCount(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "get_like_count",
		"password":     "testpasswd",
		"email":        "get_like_count@yahoo.com",
		"phone_number": "+79250000023",
		"name":         "Get",
		"surname":      "LikeCount",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "get_like_count",
		"password": "testpasswd",
	})

	post := map[string]any{
		"text": "test post with likes",
	}
	postId := createPostOk(t, id, post, jwtAuth(token))

	const LIKE_COUNT = 5
	for i := range LIKE_COUNT {
		registerUserOk(t, map[string]any{
			"login":        fmt.Sprintf("get_like_count_%d", i),
			"password":     "testpasswd",
			"email":        fmt.Sprintf("get_like_count%d@yahoo.com", i),
			"phone_number": fmt.Sprintf("+7926000000%d", i),
			"name":         "Get",
			"surname":      fmt.Sprintf("LikeCount%d", i),
		})
		likerToken := authenticateOk(t, map[string]any{
			"login":    fmt.Sprintf("get_like_count_%d", i),
			"password": "testpasswd",
		})
		newLikeOk(t, postId, jwtAuth(likerToken))
	}

	time.Sleep(5 * time.Second)

	likeCount := getMetricOk(t, postId, "like_count")
	require.Equal(t, LIKE_COUNT, likeCount)
}

func TestGetCommentCount(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "get_comment_count",
		"password":     "testpasswd",
		"email":        "get_comment_count@yahoo.com",
		"phone_number": "+79250000024",
		"name":         "Get",
		"surname":      "CommentCount",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "get_comment_count",
		"password": "testpasswd",
	})

	post := map[string]any{
		"text": "test post with comments",
	}
	postId := createPostOk(t, id, post, jwtAuth(token))

	const COMMENT_COUNT = 5
	for i := range COMMENT_COUNT {
		comment := map[string]any{
			"content": fmt.Sprintf("comment #%d", i),
		}
		newCommentOk(t, postId, comment, jwtAuth(token))
	}

	time.Sleep(5 * time.Second)

	commentCount := getMetricOk(t, postId, "comment_count")
	require.Equal(t, COMMENT_COUNT, commentCount)
}

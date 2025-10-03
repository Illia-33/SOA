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

func tryGetMetricDynamics(t *testing.T, postId int, metric string) *http.Response {
	body := map[string]any{
		"metric": metric,
	}
	return makeRequest(t, http.MethodGet, fmt.Sprintf("/post/%d/metric_dynamics", postId), body, "")
}

func getMetricDynamicsOk(t *testing.T, postId int, metric string) []map[string]any {
	resp := tryGetMetricDynamics(t, postId, metric)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	dynamicsAny := (responseBodyToMap(t, resp)["dynamics"]).([]any)

	dynamics := make([]map[string]any, len(dynamicsAny))
	for i := range dynamicsAny {
		dynamics[i] = dynamicsAny[i].(map[string]any)
	}

	return dynamics
}

func tryGetTop10Posts(t *testing.T, metric string) *http.Response {
	body := map[string]any{
		"metric": metric,
	}
	return makeRequest(t, http.MethodGet, "/top10/posts", body, "")
}

func getTop10PostsOk(t *testing.T, metric string) []map[string]any {
	resp := tryGetTop10Posts(t, metric)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	top10Any := (responseBodyToMap(t, resp)["posts"]).([]any)

	top10 := make([]map[string]any, len(top10Any))
	for i := range top10Any {
		top10[i] = top10Any[i].(map[string]any)
	}

	return top10
}

func tryGetTop10Users(t *testing.T, metric string) *http.Response {
	body := map[string]any{
		"metric": metric,
	}
	return makeRequest(t, http.MethodGet, "/top10/users", body, "")
}

func getTop10UsersOk(t *testing.T, metric string) []map[string]any {
	resp := tryGetTop10Users(t, metric)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	top10Any := (responseBodyToMap(t, resp)["users"]).([]any)

	top10 := make([]map[string]any, len(top10Any))
	for i := range top10Any {
		top10[i] = top10Any[i].(map[string]any)
	}

	return top10
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

	time.Sleep(10 * time.Second)

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

	time.Sleep(10 * time.Second)

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

	time.Sleep(10 * time.Second)

	commentCount := getMetricOk(t, postId, "comment_count")
	require.Equal(t, COMMENT_COUNT, commentCount)
}

func TestGetViewCountDynamics(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "get_view_count_dynamics",
		"password":     "testpasswd",
		"email":        "get_view_count_dynamics@yahoo.com",
		"phone_number": "+79250000025",
		"name":         "Get",
		"surname":      "ViewCountDynamics",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "get_view_count_dynamics",
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

	time.Sleep(10 * time.Second)

	dynamics := getMetricDynamicsOk(t, postId, "view_count")
	require.Equal(t, 1, len(dynamics))
	require.Equal(t, VIEW_COUNT, int(dynamics[0]["count"].(float64)))
}

func TestGetLikeCountDynamics(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "get_like_count_dynamics",
		"password":     "testpasswd",
		"email":        "get_like_count_dynamics@yahoo.com",
		"phone_number": "+79250000026",
		"name":         "Get",
		"surname":      "ViewCountDynamics",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "get_like_count_dynamics",
		"password": "testpasswd",
	})

	post := map[string]any{
		"text": "test post with likes",
	}
	postId := createPostOk(t, id, post, jwtAuth(token))

	const LIKE_COUNT = 5
	for i := range LIKE_COUNT {
		registerUserOk(t, map[string]any{
			"login":        fmt.Sprintf("get_like_count_dynamics_%d", i),
			"password":     "testpasswd",
			"email":        fmt.Sprintf("get_like_count_dynamics_%d@yahoo.com", i),
			"phone_number": fmt.Sprintf("+7927000000%d", i),
			"name":         "Get",
			"surname":      fmt.Sprintf("LikeCountDynamics%d", i),
		})
		likerToken := authenticateOk(t, map[string]any{
			"login":    fmt.Sprintf("get_like_count_dynamics_%d", i),
			"password": "testpasswd",
		})
		newLikeOk(t, postId, jwtAuth(likerToken))
	}

	time.Sleep(10 * time.Second)

	dynamics := getMetricDynamicsOk(t, postId, "like_count")
	require.Equal(t, 1, len(dynamics))
	require.Equal(t, LIKE_COUNT, int(dynamics[0]["count"].(float64)))
}

func TestGetCommentCountDynamics(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "get_comment_count_dynamics",
		"password":     "testpasswd",
		"email":        "get_comment_count_dynamics@yahoo.com",
		"phone_number": "+79250000027",
		"name":         "Get",
		"surname":      "CommentCountDynamics",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "get_comment_count_dynamics",
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

	time.Sleep(10 * time.Second)

	dynamics := getMetricDynamicsOk(t, postId, "comment_count")
	require.Equal(t, 1, len(dynamics))
	require.Equal(t, COMMENT_COUNT, int(dynamics[0]["count"].(float64)))
}

func TestGetTop10PostsByViewCount(t *testing.T) {
	const TOP_POST_VIEW_COUNT = 20
	postIds := make([]int, 10)
	for postNum := range 10 {
		id := registerUserOk(t, map[string]any{
			"login":        fmt.Sprintf("top_poster_views_%d", postNum),
			"password":     "testpasswd",
			"email":        fmt.Sprintf("top_poster_views_%d@yahoo.com", postNum),
			"phone_number": fmt.Sprintf("+7925000028%d", postNum),
			"name":         "Top",
			"surname":      fmt.Sprintf("PosterViews_%d", postNum),
		})

		token := authenticateOk(t, map[string]any{
			"login":    fmt.Sprintf("top_poster_views_%d", postNum),
			"password": "testpasswd",
		})

		post := map[string]any{
			"text": fmt.Sprintf("top post #%d", postNum),
		}
		postId := createPostOk(t, id, post, jwtAuth(token))
		postIds[postNum] = postId

		for range TOP_POST_VIEW_COUNT - postNum {
			newViewOk(t, postId, jwtAuth(token))
		}
	}

	time.Sleep(30 * time.Second)

	top10Posts := getTop10PostsOk(t, "view_count")

	require.Equal(t, 10, len(top10Posts))
	for i, id := range postIds {
		require.Equal(t, id, int(top10Posts[i]["post_id"].(float64)))
		require.Equal(t, TOP_POST_VIEW_COUNT-i, int(top10Posts[i]["value"].(float64)))
	}
}

func TestGetTop10UsersByViewCount(t *testing.T) {
	const TOP_USER_VIEW_COUNT = 73
	profileIds := make([]string, 10)
	postIds := make([]int, 10)
	for i := range profileIds {
		profileIds[i] = registerUserOk(t, map[string]any{
			"login":        fmt.Sprintf("top_user_views_%d", i),
			"password":     "testpasswd",
			"email":        fmt.Sprintf("top_user_views_%d@yahoo.com", i),
			"phone_number": fmt.Sprintf("+7928000000%d", i),
			"name":         "Top",
			"surname":      fmt.Sprintf("UserViews%d", i),
		})
	}

	for userNum := range postIds {
		token := authenticateOk(t, map[string]any{
			"login":    fmt.Sprintf("top_user_views_%d", userNum),
			"password": "testpasswd",
		})
		post := map[string]any{
			"text": fmt.Sprintf("post by top user %s: #%d", profileIds[userNum], userNum),
		}
		postIds[userNum] = createPostOk(t, profileIds[userNum], post, jwtAuth(token))
	}

	for userNum := range profileIds {
		token := authenticateOk(t, map[string]any{
			"login":    fmt.Sprintf("top_user_views_%d", userNum),
			"password": "testpasswd",
		})
		for i := range postIds {
			if i == userNum {
				for range 10 - i {
					newViewOk(t, postIds[userNum], jwtAuth(token))
				}
			} else {
				for range 7 {
					newViewOk(t, postIds[i], jwtAuth(token))
				}
			}
		}
	}

	time.Sleep(30 * time.Second)

	top10Users := getTop10UsersOk(t, "view_count")
	require.Equal(t, 10, len(top10Users))
	for i, user := range top10Users {
		require.Equal(t, profileIds[i], user["user_id"].(string))
		require.Equal(t, TOP_USER_VIEW_COUNT-i, int(user["value"].(float64)))
	}
}

package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tryGetPageSettings(t *testing.T, profileId string) *http.Response {
	return makeRequest(t, http.MethodGet, fmt.Sprintf("/profile/%s/page/settings", profileId), nil, "")
}

func getPageSettingsOk(t *testing.T, profileId string) map[string]any {
	resp := tryGetPageSettings(t, profileId)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return responseBodyToMap(t, resp)
}

func tryEditPageSettings(t *testing.T, profileId string, editRequest map[string]any, auth string) *http.Response {
	return makeRequest(t, http.MethodPut, fmt.Sprintf("/profile/%s/page/settings", profileId), editRequest, auth)
}

func editPageSettingsOk(t *testing.T, profileId string, editRequest map[string]any, authorization string) {
	resp := tryEditPageSettings(t, profileId, editRequest, authorization)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryCreatePost(t *testing.T, profileId string, createPostRequest map[string]any, auth string) *http.Response {
	return makeRequest(t, http.MethodPost, fmt.Sprintf("/profile/%s/page/posts", profileId), createPostRequest, auth)
}

func createPostOk(t *testing.T, profileId string, createPostRequest map[string]any, authorization string) int {
	resp := tryCreatePost(t, profileId, createPostRequest, authorization)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return int(responseBodyToMap(t, resp)["post_id"].(float64))
}

func tryGetPost(t *testing.T, postId int, auth string) *http.Response {
	return makeRequest(t, http.MethodGet, fmt.Sprintf("/post/%d", postId), nil, auth)
}

func getPostOk(t *testing.T, postId int, auth string) map[string]any {
	resp := tryGetPost(t, postId, auth)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return responseBodyToMap(t, resp)
}

func tryGetPosts(t *testing.T, profileId string, getPostsRequest map[string]any, auth string) *http.Response {
	return makeRequest(t, http.MethodGet, fmt.Sprintf("/profile/%s/page/posts", profileId), getPostsRequest, auth)
}

func getPostsOk(t *testing.T, profileId string, getPostsRequest map[string]any, auth string) map[string]any {
	resp := tryGetPosts(t, profileId, getPostsRequest, auth)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return responseBodyToMap(t, resp)
}

func tryEditPost(t *testing.T, postId int, editPostRequest map[string]any, auth string) *http.Response {
	return makeRequest(t, http.MethodPut, fmt.Sprintf("/post/%d", postId), editPostRequest, auth)
}

func editPostOk(t *testing.T, postId int, editPostRequest map[string]any, auth string) {
	resp := tryEditPost(t, postId, editPostRequest, auth)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryDeletePost(t *testing.T, postId int, auth string) *http.Response {
	return makeRequest(t, http.MethodDelete, fmt.Sprintf("/post/%d", postId), nil, auth)
}

func deletePostOk(t *testing.T, postId int, auth string) {
	resp := tryDeletePost(t, postId, auth)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryNewView(t *testing.T, postId int, auth string) *http.Response {
	return makeRequest(t, http.MethodPost, fmt.Sprintf("/post/%d/views", postId), nil, auth)
}

func newViewOk(t *testing.T, postId int, auth string) {
	resp := tryNewView(t, postId, auth)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryNewLike(t *testing.T, postId int, auth string) *http.Response {
	return makeRequest(t, http.MethodPost, fmt.Sprintf("/post/%d/likes", postId), nil, auth)
}

func newLikeOk(t *testing.T, postId int, auth string) {
	resp := tryNewLike(t, postId, auth)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEditPageSettings(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "edit_page_settings",
		"password":     "testpasswd",
		"email":        "edit_page_settings@yahoo.com",
		"phone_number": "+79250000011",
		"name":         "Edit",
		"surname":      "PageSettings",
	})

	token := authenticateOk(t, map[string]any{
		"login":    "edit_page_settings",
		"password": "testpasswd",
	})

	pageSettings := getPageSettingsOk(t, id)

	visibleForUnauthorized := pageSettings["visible_for_unauthorized"].(bool)
	commentsEnabled := pageSettings["comments_enabled"].(bool)
	anyoneCanPost := pageSettings["anyone_can_post"].(bool)
	checkPageSettings := func(pageSettingsResponse map[string]any) {
		assert.Equal(t, visibleForUnauthorized, pageSettingsResponse["visible_for_unauthorized"].(bool))
		assert.Equal(t, commentsEnabled, pageSettingsResponse["comments_enabled"].(bool))
		assert.Equal(t, anyoneCanPost, pageSettingsResponse["anyone_can_post"].(bool))
	}

	editPageSettingsOk(t, id, map[string]any{
		"visible_for_unauthorized": !visibleForUnauthorized,
	}, jwtAuth(token))
	visibleForUnauthorized = !visibleForUnauthorized
	resp := getPageSettingsOk(t, id)
	checkPageSettings(resp)

	editPageSettingsOk(t, id, map[string]any{
		"comments_enabled": !commentsEnabled,
	}, jwtAuth(token))
	commentsEnabled = !commentsEnabled
	resp = getPageSettingsOk(t, id)
	checkPageSettings(resp)

	editPageSettingsOk(t, id, map[string]any{
		"anyone_can_post": !anyoneCanPost,
	}, jwtAuth(token))
	anyoneCanPost = !anyoneCanPost
	resp = getPageSettingsOk(t, id)
	checkPageSettings(resp)
}

func TestCreatePost(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "create_post",
		"password":     "testpasswd",
		"email":        "create_post@yahoo.com",
		"phone_number": "+79250000012",
		"name":         "Create",
		"surname":      "Post",
	})

	token := authenticateOk(t, map[string]any{
		"login":    "create_post",
		"password": "testpasswd",
	})

	post := map[string]any{
		"text": "new test post!",
	}

	postId := createPostOk(t, id, post, jwtAuth(token))

	postResponse := getPostOk(t, postId, jwtAuth(token))
	assert.Equal(t, post["text"].(string), postResponse["text"].(string))
}

func TestCreatePostRepost(t *testing.T) {
	idPoster := registerUserOk(t, map[string]any{
		"login":        "create_post_repost1",
		"password":     "testpasswd",
		"email":        "create_post_repost1@yahoo.com",
		"phone_number": "+79250000013",
		"name":         "Create",
		"surname":      "PostRepost1",
	})
	tokenPoster := authenticateOk(t, map[string]any{
		"login":    "create_post_repost1",
		"password": "testpasswd",
	})

	idReposter := registerUserOk(t, map[string]any{
		"login":        "create_post_repost2",
		"password":     "testpasswd",
		"email":        "create_post_repost2@yahoo.com",
		"phone_number": "+79250000014",
		"name":         "Create",
		"surname":      "PostRepost2",
	})
	tokenReposter := authenticateOk(t, map[string]any{
		"login":    "create_post_repost2",
		"password": "testpasswd",
	})

	post := map[string]any{
		"text": "new test post!",
	}
	postId := createPostOk(t, idPoster, post, jwtAuth(tokenPoster))

	repost := map[string]any{
		"repost": postId,
	}
	repostId := createPostOk(t, idReposter, repost, jwtAuth(tokenReposter))

	postResponse := getPostOk(t, repostId, jwtAuth(tokenReposter))
	assert.Equal(t, postId, int(postResponse["source_post_id"].(float64)))
}

func TestGetPosts(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "get_posts",
		"password":     "testpasswd",
		"email":        "get_posts@yahoo.com",
		"phone_number": "+79250000015",
		"name":         "Get",
		"surname":      "Posts",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "get_posts",
		"password": "testpasswd",
	})

	const post_count = 37
	publishedPosts := make([]map[string]any, 0, post_count)
	for i := range post_count {
		post := map[string]any{
			"text": fmt.Sprintf("post number %d", i),
		}

		postId := createPostOk(t, id, post, jwtAuth(token))
		post["id"] = postId

		publishedPosts = append(publishedPosts, post)
	}

	idx := len(publishedPosts) - 1
	pageToken := ""
	for {
		getPostsResponse := getPostsOk(t, id, map[string]any{
			"page_token": pageToken,
		}, jwtAuth(token))

		postsResponse := getPostsResponse["posts"].([]any)
		for _, postAny := range postsResponse {
			post := postAny.(map[string]any)
			require.GreaterOrEqual(t, idx, 0)
			assert.Equal(t, publishedPosts[idx]["id"].(int), int(post["id"].(float64)))
			assert.Equal(t, publishedPosts[idx]["text"].(string), post["text"].(string))
			idx--
		}

		pageToken = getPostsResponse["next_page_token"].(string)
		if pageToken == "" {
			break
		}
	}
}

func TestEditPost(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "edit_post",
		"password":     "testpasswd",
		"email":        "edit_post@yahoo.com",
		"phone_number": "+79250000016",
		"name":         "Edit",
		"surname":      "Post",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "edit_post",
		"password": "testpasswd",
	})

	postId := createPostOk(t, id, map[string]any{"text": "before edit"}, jwtAuth(token))

	editPostOk(t, postId, map[string]any{"text": "after edit"}, jwtAuth(token))
	editedPost := getPostOk(t, postId, jwtAuth(token))
	require.Equal(t, "after edit", editedPost["text"].(string))
}

func TestDeletePost(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "delete_post",
		"password":     "testpasswd",
		"email":        "delete_post@yahoo.com",
		"phone_number": "+79250000017",
		"name":         "Delete",
		"surname":      "Post",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "delete_post",
		"password": "testpasswd",
	})

	postId := createPostOk(t, id, map[string]any{"text": "post content"}, jwtAuth(token))

	deletePostOk(t, postId, jwtAuth(token))
	deleteResponse := tryGetPost(t, postId, jwtAuth(token))
	require.Equal(t, http.StatusNotFound, deleteResponse.StatusCode)
}

func TestNewView(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "new_view",
		"password":     "testpasswd",
		"email":        "new_view@yahoo.com",
		"phone_number": "+79250000018",
		"name":         "New",
		"surname":      "View",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "new_view",
		"password": "testpasswd",
	})

	postId := createPostOk(t, id, map[string]any{"text": "post content"}, jwtAuth(token))
	newViewOk(t, postId, jwtAuth(token))
	post := getPostOk(t, postId, jwtAuth(token))
	require.Equal(t, 1, int(post["views_count"].(float64)))
}

func TestNewLike(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "new_like",
		"password":     "testpasswd",
		"email":        "new_like@yahoo.com",
		"phone_number": "+79250000019",
		"name":         "New",
		"surname":      "Like",
	})
	token := authenticateOk(t, map[string]any{
		"login":    "new_like",
		"password": "testpasswd",
	})

	postId := createPostOk(t, id, map[string]any{"text": "post content"}, jwtAuth(token))
	newLikeOk(t, postId, jwtAuth(token))
}

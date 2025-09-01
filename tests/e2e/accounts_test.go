package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tryRegisterUser(t *testing.T, registerRequest map[string]any) *http.Response {
	return makeRequest(t, http.MethodPost, "/profile", registerRequest, "")
}

func registerUserOk(t *testing.T, registerRequest map[string]any) string {
	resp := tryRegisterUser(t, registerRequest)
	return responseBodyToMap(t, resp)["profile_id"].(string)
}

func tryGetProfileInfo(t *testing.T, profileId string) *http.Response {
	return makeRequest(t, http.MethodGet, fmt.Sprintf("/profile/%s", profileId), nil, "")
}

func getProfileInfoOk(t *testing.T, profileId string) map[string]any {
	resp := tryGetProfileInfo(t, profileId)
	return responseBodyToMap(t, resp)
}

func tryEditProfile(t *testing.T, profileId string, editRequest map[string]any, auth string) *http.Response {
	return makeRequest(t, http.MethodPut, fmt.Sprintf("/profile/%s", profileId), editRequest, auth)
}

func editProfileOk(t *testing.T, profileId string, editRequest map[string]any, authorization string) {
	resp := tryEditProfile(t, profileId, editRequest, authorization)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryDeleteProfile(t *testing.T, profileId string, auth string) *http.Response {
	return makeRequest(t, http.MethodDelete, fmt.Sprintf("/profile/%s", profileId), nil, auth)
}

func deleteProfileOk(t *testing.T, profileId string, authorization string) {
	resp := tryDeleteProfile(t, profileId, authorization)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryCreateApiToken(t *testing.T, createTokenRequest map[string]any) *http.Response {
	return makeRequest(t, http.MethodPost, "/api_token", createTokenRequest, "")
}

func createApiTokenOk(t *testing.T, createTokenRequest map[string]any) string {
	resp := tryCreateApiToken(t, createTokenRequest)
	return responseBodyToMap(t, resp)["token"].(string)
}

func TestRegister(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "register_test",
		"password":     "testpasswd",
		"email":        "register_test@yahoo.com",
		"phone_number": "+79250000000",
		"name":         "Register",
		"surname":      "Test",
	})

	profileInfo := getProfileInfoOk(t, id)
	assert.Equal(t, "Register", profileInfo["name"].(string))
	assert.Equal(t, "Test", profileInfo["surname"].(string))
}

func TestAuthSimpleLogin(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "auth_simple_login",
		"password":     "testpasswd",
		"email":        "auth_simple_login@yahoo.com",
		"phone_number": "+79250000001",
		"name":         "Auth",
		"surname":      "SimpleLogin",
	})

	token := authenticateOk(t, map[string]any{
		"login":    "auth_simple_login",
		"password": "testpasswd",
	})

	editProfileOk(t, id, map[string]any{
		"bio": "new bio",
	}, jwtAuth(token))
}

func TestAuthSimplePhoneNumber(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "auth_simple_phone_number",
		"password":     "testpasswd",
		"email":        "auth_simple_phone_number@yahoo.com",
		"phone_number": "+79250000002",
		"name":         "Auth",
		"surname":      "SimplePhoneNumber",
	})

	token := authenticateOk(t, map[string]any{
		"phone_number": "+79250000002",
		"password":     "testpasswd",
	})

	editProfileOk(t, id, map[string]any{
		"bio": "new bio",
	}, jwtAuth(token))
}

func TestAuthSimpleEmail(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "auth_simple_email",
		"password":     "testpasswd",
		"email":        "auth_simple_email@yahoo.com",
		"phone_number": "+79250000003",
		"name":         "Auth",
		"surname":      "SimpleEmail",
	})

	token := authenticateOk(t, map[string]any{
		"email":    "auth_simple_email@yahoo.com",
		"password": "testpasswd",
	})

	editProfileOk(t, id, map[string]any{
		"bio": "new bio",
	}, jwtAuth(token))
}

// func TestAuthJwtTimeout(t *testing.T) {
// 	id := registerUserOk(t, map[string]any{
// 		"login":        "auth_jwt_timeout",
// 		"password":     "testpasswd",
// 		"email":        "auth_jwt_timeout@yahoo.com",
// 		"phone_number": "+79250000004",
// 		"name":         "Auth",
// 		"surname":      "JwtTimeout",
// 	})

// 	token := authenticateOk(t, map[string]any{
// 		"email":    "auth_jwt_timeout@yahoo.com",
// 		"password": "testpasswd",
// 	})

// 	time.Sleep(30 * time.Second)

// 	resp := tryEditProfile(t, id, map[string]any{
// 		"bio": "new bio",
// 	}, jwtAuth(token))
// 	require.Equal(t, http.StatusForbidden, resp.StatusCode)
// }

func TestAuthWrongToken(t *testing.T) {
	id1 := registerUserOk(t, map[string]any{
		"login":        "auth_wrong_token_1",
		"password":     "testpasswd",
		"email":        "auth_wrong_token_1@yahoo.com",
		"phone_number": "+79250000005",
		"name":         "Auth",
		"surname":      "WrongToken1",
	})

	registerUserOk(t, map[string]any{
		"login":        "auth_wrong_token_2",
		"password":     "testpasswd",
		"email":        "auth_wrong_token_2@yahoo.com",
		"phone_number": "+79250000006",
		"name":         "Auth",
		"surname":      "WrongToken2",
	})

	token := authenticateOk(t, map[string]any{
		"email":    "auth_wrong_token_2@yahoo.com",
		"password": "testpasswd",
	})

	resp := tryEditProfile(t, id1, map[string]any{
		"bio": "new bio",
	}, jwtAuth(token))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestUpdate(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "update",
		"password":     "testpasswd",
		"email":        "update@yahoo.com",
		"phone_number": "+79250000007",
		"name":         "Test",
		"surname":      "Update",
	})

	token := authenticateOk(t, map[string]any{
		"login":    "update",
		"password": "testpasswd",
	})

	editProfileOk(t, id, map[string]any{
		"bio": "new bio",
	}, jwtAuth(token))

	profileInfo := getProfileInfoOk(t, id)
	assert.Equal(t, "new bio", profileInfo["bio"].(string))
}

func TestDelete(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "delete",
		"password":     "testpasswd",
		"email":        "delete@yahoo.com",
		"phone_number": "+79250000008",
		"name":         "Test",
		"surname":      "Delete",
	})

	token := authenticateOk(t, map[string]any{
		"login":    "delete",
		"password": "testpasswd",
	})

	deleteProfileOk(t, id, jwtAuth(token))

	profileInfo := tryGetProfileInfo(t, id)
	assert.Equal(t, http.StatusNotFound, profileInfo.StatusCode)
}

func TestApiTokenSimple(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "api_token_simple",
		"password":     "testpasswd",
		"email":        "api_token_simple@yahoo.com",
		"phone_number": "+79250000009",
		"name":         "Test",
		"surname":      "ApiTokenSimple",
	})

	token := createApiTokenOk(t, map[string]any{
		"auth": map[string]any{
			"login":    "api_token_simple",
			"password": "testpasswd",
		},
		"read_access":  true,
		"write_access": true,
		"ttl":          "1h",
	})

	editProfileOk(t, id, map[string]any{
		"bio": "new bio",
	}, soaTokenAuth(token))
}

// func TestApiTokenTimeout(t *testing.T) {
// 	id := registerUserOk(t, map[string]any{
// 		"login":        "api_token_timeout",
// 		"password":     "testpasswd",
// 		"email":        "api_token_timeout@yahoo.com",
// 		"phone_number": "+79250000010",
// 		"name":         "Test",
// 		"surname":      "ApiTokenTimeout",
// 	})

// 	token := createApiTokenOk(t, map[string]any{
// 		"auth": map[string]any{
// 			"login":    "api_token_timeout",
// 			"password": "testpasswd",
// 		},
// 		"read_access":  true,
// 		"write_access": true,
// 		"ttl":          "5s",
// 	})

// 	time.Sleep(6 * time.Second)

// 	resp := tryEditProfile(t, id, map[string]any{
// 		"bio": "new bio",
// 	}, soaTokenAuth(token))
// 	require.Equal(t, http.StatusForbidden, resp.StatusCode)
// }

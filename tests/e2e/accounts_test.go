package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gateway_api_url = "http://localhost:8080/api/v1"

func makeReader(request any) io.Reader {
	b, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer([]byte(b))
}

func parseResponse(response *http.Response) (map[string]any, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var unmarshalled map[string]any
	err = json.Unmarshal(body, &unmarshalled)
	if err != nil {
		return nil, err
	}

	return unmarshalled, nil
}

func jwtAuth(jwt string) string {
	return fmt.Sprintf("Bearer %s", jwt)
}

func soaTokenAuth(soa string) string {
	return fmt.Sprintf("SoaToken %s", soa)
}

func tryRegisterUser(t *testing.T, registerRequest map[string]any) *http.Response {
	resp, err := http.Post(fmt.Sprintf("%s/profile", gateway_api_url), "application/json", makeReader(registerRequest))
	require.NoError(t, err, "error while sending register request")
	return resp
}

func registerUserOk(t *testing.T, registerRequest map[string]any) string {
	resp := tryRegisterUser(t, registerRequest)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	registerResponse, err := parseResponse(resp)
	require.NoError(t, err, "error while unparsing response")

	profileId := registerResponse["profile_id"].(string)
	return profileId
}

func tryGetProfileInfo(t *testing.T, profileId string) *http.Response {
	resp, err := http.Get(fmt.Sprintf("%s/profile/%s", gateway_api_url, profileId))
	require.NoError(t, err, "error while sending get profile request")
	return resp
}

func getProfileInfoOk(t *testing.T, profileId string) map[string]any {
	resp := tryGetProfileInfo(t, profileId)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	profileResponse, err := parseResponse(resp)
	require.NoError(t, err, "error while unparsing response")

	return profileResponse
}

func tryAuthenticate(t *testing.T, authRequest map[string]any) *http.Response {
	resp, err := http.Post(fmt.Sprintf("%s/auth", gateway_api_url), "application/json", makeReader(authRequest))
	require.NoError(t, err, "error while sending register request")
	return resp
}

func authenticateOk(t *testing.T, authRequest map[string]any) string {
	resp := tryAuthenticate(t, authRequest)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	authResponse, err := parseResponse(resp)
	require.NoError(t, err, "error while unparsing response")

	return authResponse["token"].(string)
}

func tryEditProfile(t *testing.T, profileId string, editRequest map[string]any, authorization string) *http.Response {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/profile/%s", gateway_api_url, profileId), makeReader(editRequest))
	require.NoError(t, err, "error while creating request")
	req.Header.Add("Content-Type", "application/json")
	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "error while sending register request")
	return resp
}

func editProfileOk(t *testing.T, profileId string, editRequest map[string]any, authorization string) {
	resp := tryEditProfile(t, profileId, editRequest, authorization)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryDeleteProfile(t *testing.T, profileId string, authorization string) *http.Response {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/profile/%s", gateway_api_url, profileId), nil)
	require.NoError(t, err, "error while creating request")
	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "error while sending register request")
	return resp
}

func deleteProfileOk(t *testing.T, profileId string, authorization string) {
	resp := tryDeleteProfile(t, profileId, authorization)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func tryCreateApiToken(t *testing.T, createTokenRequest map[string]any) *http.Response {
	resp, err := http.Post(fmt.Sprintf("%s/api_token", gateway_api_url), "application/json", makeReader(createTokenRequest))
	require.NoError(t, err, "error while creating request")
	return resp
}

func createApiTokenOk(t *testing.T, createTokenRequest map[string]any) string {
	resp := tryCreateApiToken(t, createTokenRequest)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	createTokenResponse, err := parseResponse(resp)
	require.NoError(t, err, "error while parsing response")
	return createTokenResponse["token"].(string)
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

func TestAuthJwtTimeout(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "auth_jwt_timeout",
		"password":     "testpasswd",
		"email":        "auth_jwt_timeout@yahoo.com",
		"phone_number": "+79250000004",
		"name":         "Auth",
		"surname":      "JwtTimeout",
	})

	token := authenticateOk(t, map[string]any{
		"email":    "auth_jwt_timeout@yahoo.com",
		"password": "testpasswd",
	})

	time.Sleep(30 * time.Second)

	resp := tryEditProfile(t, id, map[string]any{
		"bio": "new bio",
	}, jwtAuth(token))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

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

func TestApiTokenTimeout(t *testing.T) {
	id := registerUserOk(t, map[string]any{
		"login":        "api_token_timeout",
		"password":     "testpasswd",
		"email":        "api_token_timeout@yahoo.com",
		"phone_number": "+79250000010",
		"name":         "Test",
		"surname":      "ApiTokenTimeout",
	})

	token := createApiTokenOk(t, map[string]any{
		"auth": map[string]any{
			"login":    "api_token_timeout",
			"password": "testpasswd",
		},
		"read_access":  true,
		"write_access": true,
		"ttl":          "5s",
	})

	time.Sleep(6 * time.Second)

	resp := tryEditProfile(t, id, map[string]any{
		"bio": "new bio",
	}, soaTokenAuth(token))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

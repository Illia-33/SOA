package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var gateway_api_url = "http://localhost:8080/api/v1"

func gatewayApiUrl() string {
	return gateway_api_url
}

func makeReader(request any) io.Reader {
	if request == nil {
		return nil
	}
	b, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer([]byte(b))
}

func makeRequest(t *testing.T, method string, resourcePath string, body map[string]any, auth string) *http.Response {
	req, err := http.NewRequest(method, gatewayApiUrl()+resourcePath, makeReader(body))
	require.NoError(t, err, "error while creating request")

	if auth != "" {
		req.Header.Add("Authorization", auth)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "error while sending %s request to %s", method, resourcePath)
	return resp
}

func responseBodyToMap(t *testing.T, response *http.Response) map[string]any {
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "error while reading response body")

	var unmarshalled map[string]any
	err = json.Unmarshal(body, &unmarshalled)
	require.NoError(t, err, "error while unmarshalling json response")
	return unmarshalled
}

func jwtAuth(jwt string) string {
	return fmt.Sprintf("Bearer %s", jwt)
}

func soaTokenAuth(soa string) string {
	return fmt.Sprintf("SoaToken %s", soa)
}

func tryAuthenticate(t *testing.T, authRequest map[string]any) *http.Response {
	resp, err := http.Post(fmt.Sprintf("%s/auth", gatewayApiUrl()), "application/json", makeReader(authRequest))
	require.NoError(t, err, "error while sending register request")
	return resp
}

func authenticateOk(t *testing.T, authRequest map[string]any) string {
	resp := tryAuthenticate(t, authRequest)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	authResponse := responseBodyToMap(t, resp)
	return authResponse["token"].(string)
}

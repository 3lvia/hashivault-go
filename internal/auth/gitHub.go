package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func authGitHub(vaultAddr, githubToken string, client *http.Client) (AuthenticationResponse, error) {
	path := "auth/github/login"
	requestBody, err := githubLogin(githubToken)
	if err != nil {
		return nil, err
	}

	req, err := authReq(vaultAddr, path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("while building http request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while sending http request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading response body: %w", err)
	}
	defer resp.Body.Close()

	var response authenticationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("while unmarshalling response body: %w", err)
	}

	return response, nil
}

// githubLogin handles converting the github token string to a bytes buffer
func githubLogin(login string) (*bytes.Buffer, error) {
	return loginBuffer(&gitToken{
		Token: login,
	})
}

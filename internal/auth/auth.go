package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Authenticate(addr string, opts ...Option) (AuthenticationResponse, error) {
	collector := &optionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	client := collector.client
	if client == nil {
		client = &http.Client{}
	}

	if collector.gitHubToken != "" {
		return authGitHub(addr, collector.gitHubToken)
	}

	return nil, errors.New("no authentication method provided")
}

func authGitHub(vaultAddr, githubToken string) (AuthenticationResponse, error) {
	req, err := authReq(vaultAddr, githubToken)
	if err != nil {
		return nil, fmt.Errorf("while building http request: %w", err)
	}

	client := &http.Client{}
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

// authReq returns a http request for authenticating to Vault
func authReq(addr, ghToken string) (*http.Request, error) {
	path := "auth/github/login"
	body, err := githubLogin(ghToken)

	if err != nil {
		return nil, fmt.Errorf("while converting token to buffer: %w", err)
	}

	url := makeURL(addr, path)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("while building http request: %w", err)
	}

	return req, nil
}

// makeURL returns a correctly formatted url for Vault http requests
func makeURL(address, path string) string {
	return address + "/v1/" + path
}

// githubLogin handles converting the github token string to a bytes buffer
func githubLogin(login string) (*bytes.Buffer, error) {
	return loginBuffer(&gitToken{
		Token: login,
	})
}

// loginBuffer converts a login token to a bytes buffer
func loginBuffer(lt interface{}) (*bytes.Buffer, error) {
	js, err := json.Marshal(lt)
	if err != nil {
		return nil, fmt.Errorf("while marshaling token: %w", err)
	}

	return bytes.NewBuffer(js), nil
}

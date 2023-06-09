package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
		return authGitHub(addr, collector.gitHubToken, client)
	}

	if collector.k8sServicePath != "" && collector.k8sRole != "" {
		return authK8s(addr, collector.k8sServicePath, collector.k8sRole, client)
	}

	return nil, errors.New("no authentication method provided")
}

func authK8s(vaultAddr, k8ServicePath, role string, client *http.Client) (AuthenticationResponse, error) {
	path := "auth/" + k8ServicePath + "/login"
	requestBody, err := k8sLogin(k8ServicePath, role)
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

// authReq returns a http request for authenticating to Vault
func authReq(addr, path string, body *bytes.Buffer) (*http.Request, error) {
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

// k8sLogin handles converting the service path and role to a bytes buffer
func k8sLogin(k8ServicePath string, role string) (*bytes.Buffer, error) {
	jwt, err := getJWT(k8ServicePath)
	if err != nil {
		return nil, err
	}

	return loginBuffer(&k8sToken{
		JWT:  jwt,
		Role: role,
	})
}

// getJWT reads JSON web token from file at the service path
func getJWT(k8ServicePath string) (string, error) {
	b, err := ioutil.ReadFile(k8ServicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read jwt token from %s: %w", k8ServicePath, err)
	}

	return string(bytes.TrimSpace(b)), nil
}

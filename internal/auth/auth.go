package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Authenticate(addr string, method Method, opts ...Option) (AuthenticationResponse, error) {
	collector := &optionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	client := collector.client
	if client == nil {
		client = &http.Client{}
	}

	switch method {
	case MethodOICD:
		return authOICD(addr)
	case MethodGitHub:
		if collector.gitHubToken == "" {
			return nil, errors.New("no GitHub token provided")
		}
		return authGitHub(addr, collector.gitHubToken, client)
	case MethodK8s:
		if collector.k8sServicePath == "" || collector.k8sRole == "" {
			return nil, errors.New("no k8s service path or role provided")
		}
		return authK8s(addr, collector.k8sServicePath, collector.k8sRole, client)
	}

	return nil, errors.New("no authentication method provided")
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

// loginBuffer converts a login token to a bytes buffer
func loginBuffer(lt interface{}) (*bytes.Buffer, error) {
	js, err := json.Marshal(lt)
	if err != nil {
		return nil, fmt.Errorf("while marshaling token: %w", err)
	}

	return bytes.NewBuffer(js), nil
}

// getJWT reads JSON web token from file at the service path
func getJWT(k8ServicePath string) (string, error) {
	b, err := ioutil.ReadFile(k8ServicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read jwt token from %s: %w", k8ServicePath, err)
	}

	return string(bytes.TrimSpace(b)), nil
}

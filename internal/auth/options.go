package auth

import "net/http"

type optionsCollector struct {
	client *http.Client

	gitHubToken string

	k8sServicePath string
	k8sRole        string
}

// Option expl
type Option func(*optionsCollector)

// WithClient sets the http client to use for requests
func WithClient(client *http.Client) Option {
	return func(o *optionsCollector) {
		o.client = client
	}
}

// WithGitHubToken sets the GitHub token to use for authentication
func WithGitHubToken(token string) Option {
	return func(o *optionsCollector) {
		o.gitHubToken = token
	}
}

// WithK8s sets the Kubernetes service path and role to use for authentication
func WithK8s(servicePath, role string) Option {
	return func(o *optionsCollector) {
		o.k8sServicePath = servicePath
		o.k8sRole = role
	}
}

package auth

import "net/http"

type optionsCollector struct {
	client *http.Client

	gitHubToken string
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

package hashivault

import "net/http"

type optionsCollector struct {
	client *http.Client
}

// Option is a function that can be used to configure this package.
type Option func(*optionsCollector)

// WithClient sets the http client to use when making requests to Vault. This is useful for testing.
func WithClient(client *http.Client) Option {
	return func(o *optionsCollector) {
		o.client = client
	}
}

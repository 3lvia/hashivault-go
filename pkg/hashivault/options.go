package hashivault

import "net/http"

type optionsCollector struct {
	client *http.Client
}

type Option func(*optionsCollector)

func WithClient(client *http.Client) Option {
	return func(o *optionsCollector) {
		o.client = client
	}
}

package hashivault

import "net/http"

type optionsCollector struct {
	client *http.Client
}

type Option func(*optionsCollector)

func WithHTTPClient(client *http.Client) Option {
	return func(collector *optionsCollector) {
		collector.client = client
	}
}

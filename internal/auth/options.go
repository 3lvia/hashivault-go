package auth

import (
	"log"
	"net/http"
)

type optionsCollector struct {
	client *http.Client

	gitHubToken string

	k8sServicePath string
	k8sRole        string

	l              *log.Logger
	otelTracerName string
	fileReader     FileReaderFunc
}

// Option is a function that provides ab optional configuration for this package.
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

// WithLogger sets the logger to use for logging. If not set, a noop logger is used.
func WithLogger(l *log.Logger) Option {
	return func(o *optionsCollector) {
		o.l = l
	}
}

// WithOtelTracerName sets the name of the OpenTelemetry tracer to use for tracing. If not set, the default tracer is
// used ("go.opentelemetry.io/otel").
func WithOtelTracerName(name string) Option {
	return func(o *optionsCollector) {
		o.otelTracerName = name
	}
}

// WithFileReader sets the file reader to use for reading files. The file in question is the Kubernetes service account
// token file when authenticating using the workflow for Kubernetes. If not set, the default file reader is used. This
// option is mainly intended for testing.
func WithFileReader(reader FileReaderFunc) Option {
	return func(o *optionsCollector) {
		o.fileReader = reader
	}
}

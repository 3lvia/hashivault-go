package hashivault

import (
	"fmt"
	"net/http"
)

// New returns a new SecretsManager and also a channel that will send errors that may arise in the concurrent internal
// goroutines that will run in the whole lifetime of the service after this function. The returned error indicates that
// something went wrong during initialization, and the service will not be able to run (if it is not nil).
func New(opts ...Option) (SecretsManager, <-chan error, error) {
	c := &optionsCollector{}
	for _, opt := range opts {
		opt(c)
	}

	c.initialize()
	if err := c.validate(); err != nil {
		return nil, nil, fmt.Errorf("invalid options: %w", err)
	}

	errChan := make(chan error)
	client := c.client
	if client == nil {
		client = &http.Client{}
	}

	tokenGetter := startTokenJob(c.vaultAddress, c.vaultToken, c.gitHubToken, c.k8sMountPath, c.k8sRole, errChan, client)

	m := newManager(c.vaultAddress, tokenGetter, errChan)
	return m, errChan, nil
}

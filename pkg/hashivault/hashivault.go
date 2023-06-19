package hashivault

import (
	"fmt"
	"net/http"
	"sync"
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

	tokenGetter := func() string {
		return c.vaultToken
	}
	if c.vaultToken == "" {
		// initializedChan is used to signal that the tokenGetter has been initialized. This ensures that secrets are not
		// requested before we have a valid token. This channel should only be used once, and no actual message will ever
		// be sent on it. Instead, it will be closed when the tokenGetter has been initialized.
		initializedChan := make(chan struct{})

		tokenGetter = startTokenJob(c, errChan, initializedChan, client)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(ch <-chan struct{}, w *sync.WaitGroup) {
			defer w.Done()
			<-ch
		}(initializedChan, wg)

		wg.Wait()
	}

	m := newManager(c.vaultAddress, tokenGetter, errChan)
	return m, errChan, nil
}

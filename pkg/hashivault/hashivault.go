package hashivault

import (
	"fmt"
	"net/http"
	"os"
)

// New returns a new SecretsManager and also a channel that will send errors that may arise in the concurrent internal
// goroutines that will run in the whole lifetime of the service after this function. The returned error indicates that
// something went wrong during initialization, and the service will not be able to run (if it is not nil).
func New(opts ...Option) (SecretsManager, <-chan error, error) {
	collector := &optionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	vaultAddress, gitHubToken, k8sMountPath, k8sRole, err := initValues()
	if err != nil {
		return nil, nil, err
	}

	errChan := make(chan error)
	client := collector.client
	if client == nil {
		client = &http.Client{}
	}

	tokenGetter := startTokenJob(vaultAddress, gitHubToken, k8sMountPath, k8sRole, errChan, client)

	m := newManager(vaultAddress, tokenGetter, errChan)
	return m, errChan, nil
}

func initValues() (vaultAddress, gitHubToken, k8sMountPath, k8sRole string, err error) {
	vaultAddress = os.Getenv("VAULT_ADDR")
	gitHubToken = os.Getenv("GITHUB_TOKEN")
	k8sMountPath = os.Getenv("MOUNT_PATH")
	k8sRole = os.Getenv("ROLE")

	if vaultAddress == "" {
		err = fmt.Errorf("VAULT_ADDR not set")
	}
	if gitHubToken == "" && k8sMountPath == "" {
		err = fmt.Errorf("GITHUB_TOKEN or MOUNT_PATH not set")
	}

	return
}

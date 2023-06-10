package hashivault

import (
	"net/http"
	"os"
)

func New(opts ...Option) (SecretsManager, <-chan error, error) {
	collector := &optionsCollector{}
	for _, opt := range opts {
		opt(collector)
	}

	client := collector.client
	if client == nil {
		client = &http.Client{}
	}

	vaultAddress, gitHubToken, k8sMountPath, k8sRole := initValues()
	errChan := make(chan error)

	tokenGetter := startTokenJob(vaultAddress, gitHubToken, k8sMountPath, k8sRole, errChan, client)

	m := newManager(vaultAddress, tokenGetter, errChan)
	return m, errChan, nil
}

func initValues() (vaultAddress, gitHubToken, k8sMountPath, k8sRole string) {
	vaultAddress = os.Getenv("VAULT_ADDR")
	gitHubToken = os.Getenv("GITHUB_TOKEN")
	k8sMountPath = os.Getenv("MOUNT_PATH")
	k8sRole = os.Getenv("ROLE")
	return
}

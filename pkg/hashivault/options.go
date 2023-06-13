package hashivault

import (
	"fmt"
	"net/http"
	"os"
)

type optionsCollector struct {
	client          *http.Client
	vaultAddress    string
	gitHubToken     string
	k8sMountPath    string
	k8sRole         string
	loadFromEnvVars bool
}

// Option is a function that can be used to configure this package.
type Option func(*optionsCollector)

// WithClient sets the http client to use when making requests to Vault. This is useful for testing.
func WithClient(client *http.Client) Option {
	return func(o *optionsCollector) {
		o.client = client
	}
}

// WithVaultAddress sets the address of the Vault server to use when making requests to Vault.
func WithVaultAddress(address string) Option {
	return func(o *optionsCollector) {
		o.vaultAddress = address
	}
}

// WithGitHubToken sets the GitHub token to use when authenticating to Vault.
func WithGitHubToken(token string) Option {
	return func(o *optionsCollector) {
		o.gitHubToken = token
	}
}

// WithKubernetes sets the Kubernetes mount path and role to use when authenticating to Vault.
func WithKubernetes(mountPath, role string) Option {
	return func(o *optionsCollector) {
		o.k8sMountPath = mountPath
		o.k8sRole = role
	}
}

// LoadFromEnvVars specifies that the options should be loaded from environment variables.
func LoadFromEnvVars() Option {
	return func(o *optionsCollector) {
		o.loadFromEnvVars = true
	}
}

func (c *optionsCollector) initialize() {
	if !c.loadFromEnvVars {
		return
	}
	c.vaultAddress = os.Getenv("VAULT_ADDR")
	c.gitHubToken = os.Getenv("GITHUB_TOKEN")
	c.k8sMountPath = os.Getenv("MOUNT_PATH")
	c.k8sRole = os.Getenv("ROLE")
}

func (c *optionsCollector) validate() error {
	if c.vaultAddress == "" {
		return fmt.Errorf("VAULT_ADDR not set")
	}
	if c.gitHubToken == "" && c.k8sMountPath == "" {
		return fmt.Errorf("GITHUB_TOKEN or MOUNT_PATH not set")
	}
	return nil
}

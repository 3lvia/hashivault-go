package hashivault

import (
	"fmt"
	"github.com/3lvia/hashivault-go/internal/auth"
	"net/http"
	"os"
)

type optionsCollector struct {
	client       *http.Client
	vaultAddress string
	gitHubToken  string
	k8sMountPath string
	k8sRole      string
	useOIDC      bool
	vaultToken   string
}

// Option is a function that can be used to configure this package.
type Option func(*optionsCollector)

// WithClient sets the http client to use when making requests to Vault. This is useful for testing.
func WithClient(client *http.Client) Option {
	return func(o *optionsCollector) {
		o.client = client
	}
}

// WithOIDC sets the authentication method to OIDC.
func WithOIDC() Option {
	return func(o *optionsCollector) {
		o.useOIDC = true
	}
}

// WithVaultToken sets the Vault token to use when authenticating to Vault.
func WithVaultToken(token string) Option {
	return func(o *optionsCollector) {
		o.vaultToken = token
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

func (c *optionsCollector) authMethod() auth.Method {
	if c.vaultToken != "" {
		return auth.MethodToken
	}
	if c.useOIDC {
		return auth.MethodOICD
	}
	if c.k8sMountPath != "" {
		return auth.MethodK8s
	}

	return auth.MethodGitHub
}

func (c *optionsCollector) initialize() {
	va := os.Getenv("VAULT_ADDR")
	if va != "" {
		c.vaultAddress = va
	}

	ght := os.Getenv("GITHUB_TOKEN")
	if ght != "" {
		c.gitHubToken = ght
	}

	k8sMP := os.Getenv("MOUNT_PATH")
	if k8sMP != "" {
		c.k8sMountPath = k8sMP
	}

	k8sR := os.Getenv("ROLE")
	if k8sR != "" {
		c.k8sRole = k8sR
	}

	vt := os.Getenv("VAULT_TOKEN")
	if vt != "" {
		c.vaultToken = vt
	}
}

func (c *optionsCollector) validate() error {
	if c.vaultAddress == "" {
		return fmt.Errorf("VAULT_ADDR not set")
	}
	if c.useOIDC {
		return nil
	}
	if c.gitHubToken == "" && c.k8sMountPath == "" {
		return fmt.Errorf("GITHUB_TOKEN or MOUNT_PATH not set")
	}
	return nil
}

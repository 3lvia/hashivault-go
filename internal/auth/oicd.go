package auth

import (
	"fmt"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/vault/api"
)

const (
	defaultMount          = "oidc"
	defaultListenAddress  = "localhost"
	defaultPort           = "8250"
	defaultCallbackHost   = "localhost"
	defaultCallbackMethod = "http"

	FieldCallbackHost   = "callbackhost"
	FieldCallbackMethod = "callbackmethod"
	FieldListenAddress  = "listenaddress"
	FieldPort           = "port"
	FieldCallbackPort   = "callbackport"
	FieldSkipBrowser    = "skip_browser"
	FieldAbortOnError   = "abort_on_error"
)

type OICDHandler struct {
}

func (h *OICDHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = defaultMount
	}

	listenAddress, ok := m[FieldListenAddress]
	if !ok {
		listenAddress = defaultListenAddress
	}

	port, ok := m[FieldPort]
	if !ok {
		port = defaultPort
	}

	callbackHost, ok := m[FieldCallbackHost]
	if !ok {
		callbackHost = defaultCallbackHost
	}

	callbackMethod, ok := m[FieldCallbackMethod]
	if !ok {
		callbackMethod = defaultCallbackMethod
	}

	callbackPort, ok := m[FieldCallbackPort]
	if !ok {
		callbackPort = port
	}

	role := m["role"]

	authURL, clientNonce, err := fetchAuthURL(c, role, mount, callbackPort, callbackMethod, callbackHost)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func fetchAuthURL(c *api.Client, role, mount, callbackPort string, callbackMethod string, callbackHost string) (string, string, error) {
	var authURL string

	clientNonce, err := base62.Random(20)
	if err != nil {
		return "", "", err
	}

	redirectURI := fmt.Sprintf("%s://%s:%s/oidc/callback", callbackMethod, callbackHost, callbackPort)
	data := map[string]interface{}{
		"role":         role,
		"redirect_uri": redirectURI,
		"client_nonce": clientNonce,
	}

	secret, err := c.Logical().Write(fmt.Sprintf("auth/%s/oidc/auth_url", mount), data)
	if err != nil {
		return "", "", err
	}

	if secret != nil {
		authURL = secret.Data["auth_url"].(string)
	}

	if authURL == "" {
		return "", "", fmt.Errorf("Unable to authorize role %q with redirect_uri %q. Check Vault logs for more information.", role, redirectURI)
	}

	return authURL, clientNonce, nil
}

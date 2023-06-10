package hashivault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func newManager(vaultAddress string, tokenGetter tokenGetterFunc, errChan chan<- error) *manager {
	return &manager{
		vaultAddress: vaultAddress,
		client:       &http.Client{},
		tokenGetter:  tokenGetter,
		errChan:      errChan,
	}
}

type manager struct {
	vaultAddress string
	client       *http.Client
	tokenGetter  tokenGetterFunc
	errChan      chan<- error
}

func (m *manager) GetSecret(path string) (EvergreenSecretsFunc, error) {
	sec, err := get(path, m.vaultAddress, m.tokenGetter(), m.client)
	if err != nil {
		return nil, err
	}

	if !sec.Renewable {
		ss := &staticSecret{sec: sec}
		return ss.get, nil
	}

	es := newEvergreen(path, m.vaultAddress, sec, m.client, m.errChan)
	return es.get, nil
}

func (m *manager) SetDefaultGoogleCredentials(path, key string) error {
	return nil
}

func get(path, vaultAddress, token string, client *http.Client) (*secret, error) {
	url := makeURL(vaultAddress, path)

	req, err := secretsReq(url, token)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sec secret
	err = json.Unmarshal(body, &sec)
	if err != nil {
		return nil, err
	}

	return &sec, nil
}

// makeURL returns a correctly formatted url for Vault http requests
func makeURL(address, path string) string {
	return address + "/v1/" + path
}

// secretsReq returns a http request for getting secrets from Vault
func secretsReq(url, auth string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while building http request: %w", err)
	}

	req.Header.Set("X-Vault-Token", auth)

	return req, nil
}

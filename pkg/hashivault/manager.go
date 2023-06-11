package hashivault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		return sec.GetData, nil
	}

	es := newEvergreen(path, m.vaultAddress, sec, m.tokenGetter, m.client, m.errChan)
	return es.get, nil
}

func (m *manager) SetDefaultGoogleCredentials(path, key string) error {
	s, err := m.GetSecret(path)
	if err != nil {
		return err
	}

	sm := s()
	if _, ok := sm[key]; !ok {
		return fmt.Errorf("key %s not found in secret", key)
	}
	var encoded string
	var ok bool
	if encoded, ok = sm[key].(string); !ok {
		return fmt.Errorf("key %s is not a string", key)
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	fn := "google-credentials.json"
	if err := os.WriteFile(fn, decoded, 0644); err != nil {
		return err
	}

	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", fn); err != nil {
		return err
	}

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

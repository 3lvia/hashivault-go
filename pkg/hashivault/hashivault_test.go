package hashivault

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func Test_4real(t *testing.T) {
	//sm, errChan, err := New(
	//	WithVaultAddress("https://vault.dev-elvia.io"),
	//	WithVaultToken("secret_token"))
	//
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//go func(ec <-chan error) {
	//	e := <-ec
	//	if e != nil {
	//		t.Error(e)
	//	}
	//}(errChan)
	//
	//f, err := sm.GetSecret("kunde/kv/data/appinsights/kunde")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//m := f()
	//x := m
	//_ = x
}

func TestNew_static(t *testing.T) {
	url, client, closer := startTestServer(t)
	defer closer()

	addHandler("/v1/kunde/kv/data/appinsights/kunde", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, jsonStaticSecret)
	})

	gitHubToken := "my-github-token"

	sm, errChan, err := New(WithClient(client), WithVaultAddress(url), WithGitHubToken(gitHubToken))
	if err != nil {
		t.Fatal(err)
	}

	go func(ec <-chan error) {
		e := <-ec
		if e != nil {
			t.Error(e)
		}
	}(errChan)

	eg, err := sm.GetSecret("kunde/kv/data/appinsights/kunde")
	NoErr(t, err)

	sec := eg()
	if sec["instrumentation-key"] != "my-secret-instrumentation-key" {
		t.Errorf("expected 'value', got '%s'", sec["instrumentation-key"])
	}

	if loginCount < 1 {
		t.Errorf("expected loginCount to be 1, got %d", loginCount)
	}
}

func TestNew_dynamic(t *testing.T) {
	url, client, closer := startTestServer(t)
	defer closer()

	secretCount := 0
	addHandler("/v1/kunde/kv/data/appinsights/kunde2", func(w http.ResponseWriter, r *http.Request) {
		secretCount++
		js := fmt.Sprintf(jsonDynamicSecret, fmt.Sprintf("secret-%d", secretCount))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, js)
	})

	gitHubToken := "my-github-token"

	sm, errChan, err := New(WithClient(client), WithVaultAddress(url), WithGitHubToken(gitHubToken))
	if err != nil {
		t.Fatal(err)
	}

	go func(ec <-chan error) {
		e := <-ec
		if e != nil {
			t.Error(e)
		}
	}(errChan)

	eg, err := sm.GetSecret("kunde/kv/data/appinsights/kunde2")
	NoErr(t, err)

	sec := eg()
	if sec["instrumentation-key"] != "secret-1" {
		t.Errorf("expected 'value', got '%s'", sec["instrumentation-key"])
	}

	if loginCount < 1 {
		t.Errorf("expected loginCount to be 1, got %d", loginCount)
	}

	<-time.After(2 * time.Second)

	sec = eg()
	if sec["instrumentation-key"] != "secret-2" {
		t.Errorf("expected 'value', got '%s'", sec["instrumentation-key"])
	}
}

func NoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error, got: %v", err)
	}
}

const (
	ghVaultResponseTemplate = `{
    "request_id": "d645ddd7-3b2e-f28b-0138-512d5ff301a4",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": null,
    "wrap_info": null,
    "warnings": null,
    "auth": {
        "client_token": "%s",
        "accessor": "zIC3dwCcg8foRsVTxxxtdX570X5",
        "policies": [
            "coreteam",
            "default",
            "drops",
            "drops-extra",
            "ea",
            "ea-extra",
            "hes-extensions",
            "hes-extensions-extra",
            "leveransemotor",
            "leveransemotor-extra"
        ],
        "token_policies": [
            "coreteam",
            "default",
            "drops",
            "drops-extra",
            "ea",
            "ea-extra",
            "hes-extensions",
            "hes-extensions-extra",
            "leveransemotor",
            "leveransemotor-extra"
        ],
        "metadata": {
            "org": "abc",
            "username": "sasc"
        },
			"lease_duration": 1,
        "renewable": true,
        "entity_id": "2957867d-00ea-981d-e2e9-dd41ab412212",
        "token_type": "service",
        "orphan": true,
        "mfa_requirement": null,
        "num_uses": 0
    }
}`

	jsonStaticSecret = `{
    "request_id": "4dfb1662-c462-99f0-120e-ee61cd3b099e",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "data": {
            "instrumentation-key": "my-secret-instrumentation-key"
        },
        "metadata": {
            "created_time": "2020-08-26T14:56:35.936623451Z",
            "custom_metadata": null,
            "deletion_time": "",
            "destroyed": false,
            "version": 2
        }
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`
	jsonDynamicSecret = `{
    "request_id": "4dfb1662-c462-99f0-120e-ee61cd3b099e",
    "lease_id": "",
    "renewable": true,
    "lease_duration": 1,
    "data": {
        "data": {
            "instrumentation-key": "%s"
        },
        "metadata": {
            "created_time": "2020-08-26T14:56:35.936623451Z",
            "custom_metadata": null,
            "deletion_time": "",
            "destroyed": false,
            "version": 2
        }
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`
)

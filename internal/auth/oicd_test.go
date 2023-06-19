package auth

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"net/http"
	"testing"
	"time"
)

func TestOICDHandler_Auth(t *testing.T) {
	//mountPath := "kubernetes/runtimeservice/hes-extensions/quantflow-worker-trogstad"
	//mountPath := "oicd"
	//role := "quantflow-worker-trogstad"

	go func() {
		if err := http.ListenAndServe(":8250", nil); err != nil {
			t.Fatal(err)
		}
		fmt.Errorf("server stopped")
	}()

	handler := &oicdHandler{}
	client, err := api.NewClient(&api.Config{
		Address: "https://vault.dev-elvia.io/",
	})
	if err != nil {
		t.Fatal(err)
	}

	vals := map[string]string{
		//"mount": mountPath,
		//"role": role,
	}

	s, err := handler.Auth(client, vals)
	if err != nil {
		t.Fatal(err)
	}
	_ = s

	<-time.After(60 * time.Second)

}

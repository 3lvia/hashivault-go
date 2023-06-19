package auth

import (
	"testing"
)

func TestOICDHandler_Auth(t *testing.T) {
	addr := "https://vault.dev-elvia.io/"

	resp, err := authOICD(addr)
	if err != nil {
		t.Fatal(err)
	}

	x := resp
	_ = x

}

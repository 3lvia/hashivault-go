package hashivault

import (
	"testing"
)

func TestNew(t *testing.T) {
	v, err := New()
	NoErr(t, err)
	_ = v
}

func NoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error, got: %v", err)
	}
}

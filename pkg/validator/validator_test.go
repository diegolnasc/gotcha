package validator

import "testing"

func TestAvailableProviders(t *testing.T) {
	r := IsProviderAvailable("github")

	if !r {
		t.Errorf("expected a true return")
	}
}

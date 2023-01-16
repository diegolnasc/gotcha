package validator

import (
	"github.com/diegolnasc/gotcha/cmd"
)

var availableProviders = []cmd.Provider{
	cmd.GitHub,
}

// IsProviderAvailable check if current provider is implemented.
func IsProviderAvailable(provider string) bool {
	for _, p := range availableProviders {
		if p == cmd.Provider(provider) {
			return true
		}
	}
	return false
}

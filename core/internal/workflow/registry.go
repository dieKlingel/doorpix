package workflow

import "fmt"

type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

func (r *Registry) RegisterProvider(identifier string, provider Provider) error {
	if _, exists := r.providers[identifier]; exists {
		return fmt.Errorf("%w, identifier = %s", ErrProviderAlreadyRegistered, identifier)
	}

	r.providers[identifier] = provider
	return nil
}

func (r *Registry) GetProvider(identifier string) (Provider, bool) {
	provider, exists := r.providers[identifier]
	return provider, exists
}

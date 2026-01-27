package provider

import (
	"fmt"
)

type Factory struct {
	providers map[string]PaymentProvider
}

func NewProviderFactory() *Factory {
	factory := &Factory{
		providers: make(map[string]PaymentProvider),
	}

	return factory
}

func (f *Factory) RegisterProvider(providerID string, provider PaymentProvider) {
	f.providers[providerID] = provider
}

func (f *Factory) GetProvider(providerID string) (PaymentProvider, error) {
	provider, exists := f.providers[providerID]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", providerID)
	}
	return provider, nil
}

package provider

import (
	"fmt"

	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
)

type Factory struct {
	providers map[string]PaymentProvider
}

func NewProviderFactory(cfg *config.Config) *Factory {
	factory := &Factory{
		providers: make(map[string]PaymentProvider),
	}

	if cfg.Paypal.Enabled {
		factory.providers["paypal"] = NewPaypalProvider(cfg.Paypal)
	}
	return factory
}

func (f *Factory) GetProvider(providerID string) (PaymentProvider, error) {
	provider, exists := f.providers[providerID]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", providerID)
	}
	return provider, nil
}

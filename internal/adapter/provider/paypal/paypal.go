package paypal

import (
	"context"
	"net/http"
	"strconv"

	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"github.com/omerbeden/paymentgateway/internal/pkg/httpclient"
)

type Provider struct {
	httpClient *http.Client
}

func NewProvider(cfg config.Paypal) *Provider {
	return &Provider{
		httpClient: &http.Client{},
	}
}

type PaypalRequest struct {
	intent         string
	purchase_units []struct {
		amount struct {
			currency_code string
			value         string
		}
	}
}
type PayPalResponse struct {
	id             string
	intent         string
	status         string
	purchase_units []struct {
		reference_id string
		amount       struct {
			currency_code string
			value         string
		}
		payee struct {
			email_address string
			merchant_id   string
		}
	}
	create_time string
	links       []struct {
		href   string
		rel    string
		method string
	}
}

func (p *Provider) CreatePayment(ctx context.Context, payment *entity.Payment) (*provider.CreatePaymentResult, error) {

	body := PaypalRequest{
		intent: "CAPTURE",
		purchase_units: []struct {
			amount struct {
				currency_code string
				value         string
			}
		}{
			{
				amount: struct {
					currency_code string
					value         string
				}{
					currency_code: payment.Currency,
					value:         "100.00",
				},
			},
		},
	}

	var response PayPalResponse
	if err := httpclient.MakeRequest(httpclient.RequestParam[PaypalRequest]{
		Client: p.httpClient,
		Ctx:    ctx,
		Method: http.MethodPost,
		URL:    "url",
		Body:   &body,
	}, &response); err != nil {
		return nil, err
	}

	amount, _ := strconv.ParseFloat(response.purchase_units[0].amount.value, 64)
	return &provider.CreatePaymentResult{
		ProviderPaymentID: response.id,
		Status:            entity.PaymentStatus(response.status),
		Amount:            amount,
		Currency:          response.purchase_units[0].amount.currency_code,
		PaymentURL:        response.links[1].href,
		Metadata:          map[string]string{},
	}, nil
}

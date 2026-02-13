package paypal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/metrics"
	"github.com/omerbeden/paymentgateway/internal/pkg/httpclient"
)

type Provider struct {
	httpClient *http.Client
	cfg        config.Paypal
	metrics    *metrics.Metrics
}

const (
	pathCreatePayment        = "/v2/checkout/orders"
	pathAuthz                = "/v1/oauth2/token"
	pathVerifyEventSignature = "/v1/notifications/verify-webhook-signature"
	pathCaptureOrder         = "/v2/checkout/orders/:%s/capture"
	providerID               = "paypal"
)

func NewProvider(cfg config.Paypal, metrics *metrics.Metrics) *Provider {
	return &Provider{
		httpClient: &http.Client{},
		cfg:        cfg,
		metrics:    metrics,
	}
}

func (p *Provider) CreatePayment(ctx context.Context, payment *entity.Payment) (*provider.CreatePaymentResult, error) {
	start := time.Now()
	operation := "create_payment"

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
					value:         strconv.FormatFloat(payment.Amount, 'f', -1, 64),
				},
			},
		},
	}

	// Fetch access token and set Authorization header
	token, err := p.getAccessToken(ctx)
	if err != nil {
		p.metrics.ProviderErrors.WithLabelValues(
			payment.ProviderID,
			"api_auth_error",
		).Inc()
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")
	headers.Set("Authorization", "Bearer "+token)

	var response PayPalResponse
	if err := httpclient.MakeRequest(httpclient.RequestParam[PaypalRequest]{
		Client: p.httpClient,
		Header: &headers,
		Ctx:    ctx,
		Method: http.MethodPost,
		URL:    p.cfg.BaseURL + pathCreatePayment,
		Body:   body,
	}, &response); err != nil {
		p.metrics.ProviderErrors.WithLabelValues(
			payment.ProviderID,
			"api_error",
		).Inc()
		p.metrics.ProviderRequestsTotal.WithLabelValues(
			payment.ProviderID,
			operation,
			"error",
		).Inc()
		return nil, err
	}

	duration := time.Since(start).Seconds()
	p.metrics.ProviderRequestsTotal.WithLabelValues(
		payment.ProviderID,
		operation,
		"success",
	).Inc()

	p.metrics.ProviderRequestDuration.WithLabelValues(
		payment.ProviderID,
		operation,
	).Observe(duration)

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

func (p *Provider) Capture(ctx context.Context, id string) error {
	start := time.Now()
	operation := "capture_payment"

	var paypalResponse paypalCaptureResponse

	header := http.Header{}
	header.Set("Content-Type", "application/json")

	err := httpclient.MakeRequest(httpclient.RequestParam[any]{
		Client: p.httpClient,
		Header: &header,
		Method: http.MethodPost,
		URL:    p.cfg.BaseURL + fmt.Sprintf(pathCaptureOrder, id),
		Ctx:    ctx,
	}, &paypalResponse)

	duration := time.Since(start).Seconds()
	p.metrics.ProviderRequestDuration.WithLabelValues(
		providerID,
		operation,
	).Observe(duration)

	if paypalResponse.status == "COMPLETED" {
		p.metrics.ProviderRequestsTotal.WithLabelValues(
			providerID,
			operation,
			"success",
		).Inc()
		return nil
	}

	p.metrics.ProviderRequestsTotal.WithLabelValues(
		providerID,
		operation,
		"error",
	).Inc()

	p.metrics.ProviderErrors.WithLabelValues(
		providerID,
		"api_error",
	).Inc()

	return fmt.Errorf("error while capturing payment %w", err)
}
func (p *Provider) VerifyWebhook(ctx context.Context, webhookCtx *provider.WebhookContext) error {
	operation := "verify_webhook"
	webhookID := p.cfg.WebhookID
	start := time.Now()

	body := PaypalVerifySignatureRequest{
		webhook_id:        webhookID,
		transmission_id:   webhookCtx.Headers.Get("TRANSMISSION-ID"),
		transmission_time: webhookCtx.Headers.Get("PAYPAL-TRANSMISSION-TIME"),
		cert_url:          webhookCtx.Headers.Get("PAYPAL-CERT-URL"),
		auth_algo:         webhookCtx.Headers.Get("PAYPAL-AUTH-ALGO"),
		transmission_sig:  webhookCtx.Signature,
		webhook_event:     json.RawMessage(webhookCtx.Payload),
	}

	var response struct {
		verification_status string
	}

	token, err := p.getAccessToken(context.Background())
	if err != nil {
		p.metrics.ProviderErrors.WithLabelValues(
			providerID,
			"api_auth_error",
		).Inc()
		return err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+token)

	if err := httpclient.MakeRequest(httpclient.RequestParam[PaypalVerifySignatureRequest]{
		Client: p.httpClient,
		Header: &headers,
		URL:    p.cfg.BaseURL + pathVerifyEventSignature,
		Body:   body,
	}, &response); err != nil {
		p.metrics.ProviderErrors.WithLabelValues(
			providerID,
			"api_error",
		).Inc()
		return err
	}
	duration := time.Since(start).Seconds()
	p.metrics.ProviderRequestDuration.WithLabelValues(
		providerID,
		operation,
	).Observe(duration)

	if response.verification_status == "SUCCESS" {
		p.metrics.ProviderRequestsTotal.WithLabelValues(
			providerID,
			operation,
			"success",
		).Inc()
		return nil
	}

	p.metrics.ProviderRequestsTotal.WithLabelValues(
		providerID,
		operation,
		"error",
	).Inc()
	return fmt.Errorf("paypal webhook event verification failed")

}

func (p *Provider) ParseWebhook(payload []byte) (*provider.WebhookEvent, error) {

	var webhookData PaypalWebhookEvent
	if err := json.Unmarshal(payload, &webhookData); err != nil {
		return nil, fmt.Errorf("failed to parse paypal webhook %w", err)
	}

	createTime, err := time.Parse(time.RFC3339, webhookData.create_time)
	if err != nil {
		return nil, fmt.Errorf("failed to parse paypal webhook createtime %w", err)
	}

	total, err := strconv.ParseFloat(webhookData.resource.amount.total, 64)

	event := &provider.WebhookEvent{
		ProviderID:        "paypal",
		EventType:         webhookData.event_type,
		Amount:            total,
		Currency:          webhookData.resource.amount.currency,
		CreateTime:        createTime,
		RawPayload:        string(payload),
		ProviderPaymentID: webhookData.resource.id,
	}

	switch webhookData.event_type {
	case "CHECKOUT.ORDER.APPROVED":
		event.Status = entity.PaymentStatusPending
	case "CHECKOUT.ORDER.COMPLETED":
		event.Status = entity.PaymentStatusSucceeded
	case "CHECKOUT.PAYMENT-APPROVAL.REVERSED":
		event.Status = entity.PaymentStatusFailed
	}

	return event, nil

}

type accessTokenResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppId       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

func (p *Provider) getAccessToken(ctx context.Context) (string, error) {
	url := p.cfg.BaseURL + pathAuthz
	data := "grant_type=client_credentials"

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("Accept", "application/json")

	var response accessTokenResponse
	if err := httpclient.MakeRequest(httpclient.RequestParam[string]{
		Client:       p.httpClient,
		Header:       &header,
		Method:       http.MethodPost,
		URL:          url,
		Body:         data,
		Ctx:          ctx,
		ClientID:     p.cfg.ClientID,
		ClientSecret: p.cfg.ClientSecret,
	}, &response); err != nil {
		return "", fmt.Errorf("failed to execute token request: %w", err)
	}

	return response.AccessToken, nil
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

type paypalCaptureResponse struct {
	id             string
	status         string
	purchase_units []struct {
		reference_id string
		payments     struct {
			captures []struct {
				id     string
				status string
				amount struct {
					currency_code string
					value         string
				}
				final_capture               bool
				seller_receivable_breakdown struct {
					gross_amount struct {
						currency_code string
						value         string
					}
					paypal_fee struct {
						currency_code string
						value         string
					}
					net_amount struct {
						currency_code string
						value         string
					}
				}
				links []struct {
					href   string
					rel    string
					method string
				}
				create_time string
				update_time string
			}
		}
	}
}

type PaypalVerifySignatureRequest struct {
	webhook_id        string
	transmission_id   string
	transmission_time string
	cert_url          string
	auth_algo         string
	transmission_sig  string
	webhook_event     []byte
}

type PaypalWebhookEvent struct {
	id            string
	create_time   string
	resource_type string
	event_version string
	event_type    string
	summary       string
	resource      struct {
		id          string
		create_time string
		update_time string
		state       string
		amount      struct {
			total    string
			currency string
			details  struct {
				subtotal string
			}
		}
		parent_payment string
		valid_until    string
		links          []struct {
			href   string
			rel    string
			method string
		}
	}
	links []struct {
		href   string
		rel    string
		method string
	}
}

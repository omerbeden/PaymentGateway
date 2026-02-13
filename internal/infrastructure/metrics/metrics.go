package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	HTTPRequestTotal          *prometheus.CounterVec
	HTTPRequestDuration       *prometheus.HistogramVec
	PaymentsTotal             *prometheus.CounterVec
	PaymentAmount             *prometheus.HistogramVec
	PaymentDuration           *prometheus.HistogramVec
	ProviderRequestsTotal     *prometheus.CounterVec
	ProviderRequestDuration   *prometheus.HistogramVec
	ProviderErrors            *prometheus.CounterVec
	DBQueryDuration           *prometheus.HistogramVec
	WebhooksReceived          *prometheus.CounterVec
	WebhookProcessingDuration *prometheus.HistogramVec
}

func New() *Metrics {
	return &Metrics{
		HTTPRequestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			}, []string{"method", "endpoint"},
		),
		PaymentsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payments_total",
				Help: "Total number of payments processed",
			},
			[]string{"status", "currency", "provider"},
		),
		PaymentAmount: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_amount",
				Help:    "Payment amount distribution",
				Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
			},
			[]string{"currency", "provider"},
		),
		PaymentDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_processing_duration_seconds",
				Help:    "Time to process a payment",
				Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 30},
			},
			[]string{"provider", "status"},
		),
		ProviderRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "provider_requests_total",
				Help: "Total requests to payment providers",
			},
			[]string{"provider", "operation", "status"},
		),

		ProviderRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "provider_request_duration_seconds",
				Help:    "Provider API request duration",
				Buckets: []float64{.1, .25, .5, 1, 2, 5, 10},
			},
			[]string{"provider", "operation"},
		),
		ProviderErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "provider_errors_total",
				Help: "Total provider errors",
			},
			[]string{"provider", "error_type"},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{"operation"},
		),
		WebhooksReceived: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "webhooks_received_total",
				Help: "Total webhooks received from providers",
			},
			[]string{"provider", "event_type", "status"},
		),

		WebhookProcessingDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "webhook_processing_duration_seconds",
				Help:    "Webhook processing duration",
				Buckets: []float64{.01, .05, .1, .25, .5, 1, 2},
			},
			[]string{"provider"},
		),
	}
}

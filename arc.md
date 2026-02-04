# Payment Gateway Aggregator - Clean Architecture Structure

```
payment-gateway/
├── cmd/
│   ├── api/
│   │   └── main.go                    # HTTP API server entry point
│   ├── worker/
│   │   └── main.go                    # Background worker entry point
│   └── migrator/
│       └── main.go                    # Database migration runner
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── payment.go             # Payment entity
│   │   │   ├── transaction.go         # Transaction entity
│   │   │   ├── merchant.go            # Merchant entity
│   │   │   └── provider.go            # Payment provider entity
│   │   ├── repository/
│   │   │   ├── payment_repository.go  # Payment repository interface
│   │   │   ├── transaction_repository.go
│   │   │   └── merchant_repository.go
│   │   └── service/
│   │       └── payment_service.go     # Payment service interface
│   ├── usecase/
│   │   ├── payment/
│   │   │   ├── create_payment.go      # Create payment use case
│   │   │   ├── process_payment.go     # Process payment use case
│   │   │   ├── refund_payment.go      # Refund payment use case
│   │   │   └── get_payment_status.go  # Get payment status use case
│   │   ├── webhook/
│   │   │   └── handle_webhook.go      # Handle provider webhooks
│   │   └── reconciliation/
│   │       └── reconcile_payments.go  # Reconciliation use case
│   ├── adapter/
│   │   ├── handler/
│   │   │   ├── http/
│   │   │   │   ├── payment_handler.go     # Payment HTTP handlers
│   │   │   │   ├── webhook_handler.go     # Webhook HTTP handlers
│   │   │   │   ├── merchant_handler.go    # Merchant HTTP handlers
│   │   │   │   └── middleware/
│   │   │   │       ├── auth.go            # Authentication middleware
│   │   │   │       ├── rate_limit.go      # Rate limiting middleware
│   │   │   │       └── idempotency.go     # Idempotency middleware
│   │   ├── repository/
│   │   │   ├── postgres/
│   │   │   │   ├── payment_repository.go
│   │   │   │   ├── transaction_repository.go
│   │   │   │   └── merchant_repository.go
│   │   │   └── redis/
│   │   │       ├── cache_repository.go
│   │   │       └── idempotency_repository.go
│   │   └── provider/
│   │       ├── stripe/
│   │       │   ├── client.go              # Stripe client
│   │       │   ├── payment.go             # Stripe payment implementation
│   │       ├── iyzico/
│   │       │   ├── client.go              # Iyzico client
│   │       │   ├── payment.go             # Iyzico payment implementation
│   │       └── paypal/
│   │           ├── client.go              # PayPal client
│   │           ├── payment.go             # PayPal payment implementation
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres.go                # PostgreSQL connection
│   │   │   └── migrations/
│   │   │       ├── 001_create_payments.sql
│   │   │       ├── 002_create_transactions.sql
│   │   │       └── 003_create_merchants.sql
│   │   ├── cache/
│   │   │   └── redis.go                   # Redis connection
│   │   ├── queue/
│   │   │   └── redis_queue.go             # Redis-based queue
│   │   ├── logger/
│   │   │   └── logger.go                  # Logger implementation
│   │   └── config/
│   │       └── config.go                  # Configuration loader
│   └── pkg/
│       ├── errors/
│       │   └── errors.go                  # Custom error types
│       ├── validator/
│       │   └── validator.go               # Request validation
│       ├── crypto/
│       │   └── crypto.go                  # Encryption/decryption utilities
│       └── utils/
│           ├── idempotency.go             # Idempotency key generator
│           └── retry.go                   # Retry logic utilities
├── doc/
│   ├── openapi.yaml                   # OpenAPI specification
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── configmap.yaml
├── scripts/
│   ├── migrate.sh                         # Migration script
│   └── seed.sh                            # Seed data script
├── tests/
│   ├── integration/
│   │   ├── payment_test.go
│   │   └── webhook_test.go
│   └── e2e/
│       └── payment_flow_test.go
├── docs/
│   ├── architecture.md                    # Architecture documentation
│   ├── api.md                             # API documentation
│   └── deployment.md                      # Deployment guide
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

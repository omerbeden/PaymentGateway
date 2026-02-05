# PaymentGateway 🚀

**PaymentGateway** is a sample backend service implemented in Go. It implements a simple payment flow, webhook processing, and uses common infrastructure components (Postgres, Redis, Docker).

## 🔍 Overview

- Purpose: Showcase software design, Go idioms, testing, and deployment skills
- Built with a clean separation between `handler`, `usecase`, `repository`, and `provider` layers.
- Includes an OpenAPI spec (`doc/openapi.yaml`) and dockerized dev environment (`deployments/docker/docker-compose.yml`).

## ⚙️ Features

- REST API for creating and retrieving payments
- Webhook processing
- Provider abstraction (example: PayPal provider in `internal/provider/paypal`)
- PostgreSQL persistence and migrations (`internal/infrastructure/database/migrations`)
- Redis caching (`internal/infrastructure/cache`)
- Docker & docker-compose for easy local setup
- Unit and repository tests (see `internal/repository`)

## 🧭 Project Structure

```
payment-gateway/
├── cmd/
│   ├── api/
│   │   └── main.go                    # HTTP API server entry point
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── payment.go             
│   │   │   ├── transaction.go         
│   │   │   └── provider.go            
│   │   ├── repository/
│   │   │   ├── payment_repository.go  
│   │   │   ├── webhook_event_repository.go
│   │   └── service/
│   │       └── payment_service.go     # Payment service interface
│   ├── usecase/
│   │   ├── payment/
│   │   │   ├── create_payment.go      
│   │   │   ├── process_payment.go     
│   │   │   ├── refund_payment.go      
│   │   │   └── get_payment_status.go  
│   │   ├── webhook/
│   │   │   └── handle_webhook.go     
│   │   └── reconciliation/
│   │       └── reconcile_payments.go  # Reconciliation use case
│   ├── adapter/
│   │   ├── handler/
│   │   │   ├── http/
│   │   │   │   ├── payment_handler.go     
│   │   │   │   ├── webhook_handler.go     
│   │   │   │   └── middleware/
│   │   │   │       ├── auth.go            # Authentication middleware
│   │   │   │       ├── rate_limit.go      # Rate limiting middleware
│   │   │   │       └── idempotency.go     # Idempotency middleware
│   │   ├── repository/
│   │   │   ├── postgres/
│   │   │   │   ├── payment_repository.go
│   │   │   │   ├── transaction_repository.go
│   │   │   │   └── webhook_event_repository.go
│   │   │   └── redis/
│   │   │       ├── cache_repository.go
│   │   │       └── idempotency_repository.go
│   │   └── provider/
│   │       ├── stripe/
│   │       │   ├── payment.go             # Stripe payment implementation
│   │       ├── iyzico/
│   │       │   ├── payment.go             # Iyzico payment implementation
│   │       └── paypal/
│   │           ├── payment.go             # PayPal payment implementation
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres.go                # PostgreSQL connection
│   │   │   └── migrations/
│   │   │       ├── 001_create_payments.sql
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
│       ├── httpclient/
│       │   └── client.go                  # generic http client for api cals
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

## 🚀 Quick Start

Prerequisites: Go (1.18+), Docker & docker-compose (optional but recommended).

1. Run with Go (local):

```bash
# from repo root
go run ./cmd/api
```

2. Run with Docker Compose:

```bash
cd deployments/docker
docker-compose up --build
```

3. Run tests:

```bash
go test ./... -v
```

4. API docs: open `doc/openapi.yaml` with your preferred OpenAPI viewer (Swagger UI / Redoc).

> Note: database migrations are in `internal/infrastructure/database/migrations` — apply them to your Postgres instance before running.



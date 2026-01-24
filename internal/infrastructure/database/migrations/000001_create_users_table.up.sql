CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(255) PRIMARY KEY,    
    amount DECIMAL(10, 2) NOT NULL,    
    currency VARCHAR(3) NOT NULL,    
    idempotency_key VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    expires_at TIMESTAMP,
    metadata JSONB,

    CONSTRAINT unique_idempotency UNIQUE (idempotency_key)
);


CREATE INDEX idx_payments_idempotency_key ON payments(idempotency_key);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

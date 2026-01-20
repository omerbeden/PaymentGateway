package entity

import "time"

type Provider struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	APIKey        string `json:"api_key,omitempty"`
	SecretKey     string `json:"secret_key,omitempty"`
	WebhookSecret string `json:"webhook_secret,omitempty"`
	BaseURL       string `json:"base_url,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

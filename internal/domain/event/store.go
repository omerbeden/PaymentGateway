package event

import (
	"context"
)

type Store interface {
	Append(ctx context.Context, event DomainEvent) error
}

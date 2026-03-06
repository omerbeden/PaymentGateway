package integration

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/omerbeden/paymentgateway/internal/adapter/eventstore/mongodb"
	"github.com/omerbeden/paymentgateway/internal/domain/event"
)

// Run with: go test ./tests/integration/... -v -count=1

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("connect to mongodb: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("ping mongodb: %v", err)
	}

	db := client.Database("test_payment_gateway")

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Drop(ctx)
		client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestMongoEventStore_Append(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := mongodb.NewMongoEventStore(db)
	ctx := context.Background()

	paymentID := "pay_test_001"

	evt1 := event.NewPaymentCompletedEvent(paymentID, "USD", "stripe", "test payment", 100.00)
	if err := store.Append(ctx, evt1); err != nil {
		t.Fatalf("append event 1: %v", err)
	}

}

package database

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(ctx context.Context, cfg config.Mongo) (*mongo.Database, error) {

	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetConnectTimeout(cfg.Timeout).
		SetServerSelectionTimeout(cfg.Timeout)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(cfg.Database), nil
}

func Disconnect(ctx context.Context, db *mongo.Database) error {
	return db.Client().Disconnect(ctx)
}

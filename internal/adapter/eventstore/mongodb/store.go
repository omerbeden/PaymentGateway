package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/omerbeden/paymentgateway/internal/domain/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEventStore struct {
	collection *mongo.Collection
}

type eventDocument struct {
	AggregateID string        `bson:"aggregate_id"`
	Version     int           `bson:"version"`
	Events      []eventRecord `bson:"events"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}

type eventRecord struct {
	Version    int                    `bson:"version"`
	Type       string                 `bson:"type"`
	Data       map[string]interface{} `bson:"data"`
	OccurredAt time.Time              `bson:"occurred_at"`
}

func NewMongoEventStore(db *mongo.Database) *MongoEventStore {
	collection := db.Collection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "events.occurred_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "events.type", Value: 1}},
		},
	})

	return &MongoEventStore{collection: collection}
}

func (s *MongoEventStore) Append(ctx context.Context, devent event.DomainEvent) error {
	record, err := s.serializeEvent(devent)
	if err != nil {
		return err
	}

	filter := bson.M{"aggregate_id": devent.AggregateID()}
	update := bson.M{
		"$push":        bson.M{"events": record},
		"$inc":         bson.M{"version": 1},
		"$set":         bson.M{"updated_at": time.Now().UTC()},
		"$setOnInsert": bson.M{"_id": devent.AggregateID()},
	}
	opts := options.Update().SetUpsert(true)
	_, err = s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("mongo event store: append: %w", err)
	}

	return nil

}

func (s *MongoEventStore) serializeEvent(evt event.DomainEvent) (eventRecord, error) {

	data, err := json.Marshal(evt)
	if err != nil {
		return eventRecord{}, err
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return eventRecord{}, err
	}

	return eventRecord{
		Type:       string(evt.EventType()),
		Data:       dataMap,
		OccurredAt: evt.OccurredAt(),
	}, nil
}

func (s *MongoEventStore) deserializeEvent(record eventRecord) (event.DomainEvent, error) {
	data, err := json.Marshal(record.Data)
	if err != nil {
		return nil, err
	}
	switch record.Type {
	case string(event.PaymentCompleted):
		var evt event.PaymentCompletedEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			return nil, err
		}
		return evt, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", record.Type)
	}
}

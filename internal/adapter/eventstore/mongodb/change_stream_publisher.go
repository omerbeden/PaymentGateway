package mongodb

import (
	"context"
	"fmt"

	"github.com/omerbeden/paymentgateway/internal/domain/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChangeStreamPublisher struct {
	store     *MongoEventStore
	publisher event.Publisher
}

func NewChangeStreamPublisher(store *MongoEventStore, publisher event.Publisher) *ChangeStreamPublisher {
	return &ChangeStreamPublisher{
		store:     store,
		publisher: publisher,
	}
}

// Start begins watching the event_store collection and publishing events to Kafka.
func (c *ChangeStreamPublisher) Start(ctx context.Context) error {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "operationType", Value: bson.D{{Key: "$in", Value: bson.A{"insert", "update"}}}},
		}}},
	}
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	stream, err := c.store.collection.Watch(ctx, pipeline, opts)
	if err != nil {
		return fmt.Errorf("change stream: watch: %w", err)
	}
	defer stream.Close(ctx)

	for stream.Next(ctx) {
		var changeEvent struct {
			OperationType string `bson:"operationType"`
			FullDocument  struct {
				AggregateID string        `bson:"aggregate_id"`
				Events      []eventRecord `bson:"events"`
			} `bson:"fullDocument"`
		}

		if err := stream.Decode(&changeEvent); err != nil {
			//log
			continue
		}

		for _, record := range changeEvent.FullDocument.Events {
			evt, err := c.store.deserializeEvent(record)
			if err != nil {
				//log
				continue
			}
			topic := c.getTopicForEvent(evt)
			if err := c.publisher.Publish(ctx, topic, evt); err != nil {
				//log
				// Consider retrying or sending to a DLQ here
			}
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("change stream: error: %w", err)
	}

	return nil
}
func (c *ChangeStreamPublisher) getTopicForEvent(evt event.DomainEvent) string {
	switch evt.EventType() {
	case event.PaymentCompleted:
		return event.TopicNotificationPaymentCompleted
	default:
		return "events.unknown"
	}
}

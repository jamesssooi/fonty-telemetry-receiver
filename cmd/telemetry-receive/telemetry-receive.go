package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/jamesssooi/fonty-telemetry-receiver/pkg/fontytelemetry"
)

func main() {
	// Load config
	fontytelemetry.LoadConfig()

	// Start subscriber
	client, ctx := fontytelemetry.GetClient()
	sub := client.Subscription(fontytelemetry.Config.PubsubSubscription)
	_ = sub.Receive(ctx, receiver)
}

func receiver(ctx context.Context, m *pubsub.Message) {
	log.Printf("Received message: %v", string(m.Data))
	m.Ack()
}

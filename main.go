package main

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

var (
	topic *pubsub.Topic

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)

func main() {
	ctx := context.Background()
	gcpProjectID := "ga-test-project-503ca"
	topicID := "narcissus-mirror"
	subscriptionID := "narcissus-mirror-sub6"

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't initialize Zap logger: %v", err)
	}
	defer logger.Sync()
	logger.Info("Hello from Zap logger!")
	client, err := pubsub.NewClient(ctx, gcpProjectID)
	topic = client.Topic(topicID)

	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub, err := client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
		Topic:            topic,
		AckDeadline:      20 * time.Second,
		ExpirationPolicy: 24 * time.Hour,
	})
	if err != nil {
		log.Fatalf("CreateSubscription: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		log.Fatalf("Got message: %q\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		logger.Info("sub.Receive: %w")
	}

	log.Fatalf("Created subscription: %v\n", sub)

}

// func pullMsgs(w io.Writer, projectID, subID string) error {
// 	projectID = "ga-test-project-503ca"
// 	subID = "narcissus-mirror-sub2"
// 	ctx := context.Background()
// 	client, err := pubsub.NewClient(ctx, projectID)
// 	if err != nil {
// 		return fmt.Errorf("pubsub.NewClient: %w", err)
// 	}
// 	defer client.Close()

// 	sub := client.Subscription(subID)

// 	// Receive messages for 10 seconds, which simplifies testing.
// 	// Comment this out in production, since `Receive` should
// 	// be used as a long running operation.
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	var received int32
// 	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
// 		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
// 		atomic.AddInt32(&received, 1)
// 		msg.Ack()
// 	})
// 	if err != nil {
// 		return fmt.Errorf("sub.Receive: %w", err)
// 	}
// 	fmt.Fprintf(w, "Received %d messages\n", received)

// 	return nil
//}

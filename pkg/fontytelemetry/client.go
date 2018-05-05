package fontytelemetry

import (
	"log"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// GetClient returns a Google PubSub client and context.
func GetClient() (*pubsub.Client, context.Context) {
	ctx := context.Background()
	data, err := Asset("cfg/auth.json")
	creds, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/pubsub")
	client, err := pubsub.NewClient(ctx, Config.ProjectID, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client, ctx
}

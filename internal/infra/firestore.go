package infra

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func NewFirestoreClient() (*firestore.Client, error) {
	var opts []option.ClientOption
	ctx := context.Background()
	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")

	if os.Getenv("APP_ENV") == "local" {
		opts = append(opts, option.WithCredentialsFile("dina-dev-host-project.json"))
	}

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

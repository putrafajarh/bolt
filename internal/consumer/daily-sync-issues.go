package consumer

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub/v2"
	"github.com/bytedance/sonic"
	"github.com/putrafajarh/bolt/internal/njir"
	"github.com/putrafajarh/bolt/internal/services"
	"gorm.io/gorm"
)

type DailySyncIssuesMessage struct {
	Date string `json:"date"`
}

func DailySyncIssues(pubsubClient *pubsub.Client, firestoreClient *firestore.Client, db *gorm.DB) error {
	njir, err := njir.NewNjirClient()
	if err != nil {
		return err
	}

	parentCtx := context.Background()
	subID := "daily-sync-jira-issues"
	sub := pubsubClient.Subscriber(subID)

	return sub.Receive(parentCtx, func(ctx context.Context, m *pubsub.Message) {
		var msg DailySyncIssuesMessage
		if err := sonic.Unmarshal(m.Data, &msg); err != nil {
			return
		}

		yesterday := time.Now().Add(-time.Hour * 24)

		_, err := services.NewJiraServiceImpl(services.JiraServiceOptions{
			NjirClient:      njir,
			FirestoreClient: firestoreClient,
			DB:              db,
		}).SyncDailyIssues(ctx, yesterday)
		if err != nil {
			m.Nack()
		} else {
			m.Ack()
		}
	})
}

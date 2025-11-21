package services

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/andygrunwald/go-jira"
	"github.com/putrafajarh/bolt/internal/njir"
	"gorm.io/gorm"
)

type JiraService interface {
	SyncDailyIssues(ctx context.Context, date time.Time) ([]jira.Issue, error)
}

type JiraServiceImpl struct {
	njirClient      *njir.NjirClient
	firestoreClient *firestore.Client
	db              *gorm.DB
}

type JiraServiceOptions struct {
	NjirClient      *njir.NjirClient
	FirestoreClient *firestore.Client
	DB              *gorm.DB
}

func NewJiraServiceImpl(opts JiraServiceOptions) JiraService {
	return &JiraServiceImpl{
		njirClient:      opts.NjirClient,
		firestoreClient: opts.FirestoreClient,
		db:              opts.DB,
	}
}

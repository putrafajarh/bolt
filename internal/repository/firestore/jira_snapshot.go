package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/andygrunwald/go-jira"
	"github.com/bytedance/sonic"
	"google.golang.org/api/iterator"
)

type JiraSnapshotRepository struct {
	FirestoreClient *firestore.Client
}

func NewJiraSnapshotRepository(fs *firestore.Client) *JiraSnapshotRepository {
	return &JiraSnapshotRepository{
		FirestoreClient: fs,
	}
}

func (r *JiraSnapshotRepository) StoreSnapshots(ctx context.Context, date time.Time, jql string, sourceStatus string, issues []jira.Issue) (*firestore.DocumentSnapshot, error) {
	dateID := date.Format("2006-01-02")

	snapshotRef := r.FirestoreClient.Collection("jira_snapshots").Doc(dateID)
	err := r.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if err := tx.Set(snapshotRef, map[string]any{
			"runAt":        time.Now(),
			"jql":          jql,
			"sourceStatus": sourceStatus,
			"totalIssues":  len(issues),
		}); err != nil {
			return err
		}

		for _, issue := range issues {
			issueRef := snapshotRef.Collection("issues").Doc(issue.Key)

			b, err := sonic.Marshal(issue)
			if err != nil {
				return err
			}
			var m map[string]any
			if err := sonic.Unmarshal(b, &m); err != nil {
				return err
			}

			if err := tx.Set(issueRef, m); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return snapshotRef.Get(ctx)
}

func (r *JiraSnapshotRepository) CleanupSnapshots(ctx context.Context, days int) error {
	today := time.Now().Add(-time.Hour * 24 * time.Duration(days))
	todayID := today.Format("2006-01-02")

	iter := r.FirestoreClient.Collection("jira_snapshots").Where("runAt", "<", today).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if doc.Ref.ID <= todayID {
			if _, err := doc.Ref.Delete(ctx); err != nil {
				return err
			}
		}
	}

	return nil

}

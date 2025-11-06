package controller

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/andygrunwald/go-jira"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/putrafajarh/bolt/internal/infra"
	"github.com/putrafajarh/bolt/internal/njir"
	"golang.org/x/sync/errgroup"
)

func GetIssue(c *fiber.Ctx) error {
	client, err := njir.NewClient()
	if err != nil {
		return err
	}

	yesterday := time.Now().Add(-time.Hour * 24)
	// Use slash format for JQL and dash format for Firestore doc ID
	jqlDate := yesterday.Format("2006/01/02")
	dateID := yesterday.Format("2006-01-02")

	sourceStatus := "READY FOR SIT"
	jql := fmt.Sprintf("Status CHANGED TO '%s' ON '%s' ORDER BY assignee asc , created ASC", sourceStatus, jqlDate)

	issues, _, err := client.Issue.SearchV2JQLWithContext(
		c.UserContext(),
		jql,
		&jira.SearchOptionsV2{
			Expand:       "names",
			Fields:       []string{"assignee", "project", "reporter", "status", "creator", njir.StoryPointField, "summary", "created", "updated", "sprint", "transitions"},
			MaxResults:   200,
			FieldsByKeys: true,
		},
	)
	if err != nil {
		return err
	}

	// Fetch transitions concurrently before transaction
	g, ctx := errgroup.WithContext(c.UserContext())
	// g.SetLimit(8)
	for i := range issues {
		g.Go(func() error {
			transitions, _, err := client.Issue.GetTransitionsWithContext(ctx, issues[i].ID)
			if err != nil {
				return err
			}
			issues[i].Transitions = transitions
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	// Persist snapshot and issues to Firestore
	fs, err := infra.NewFirestoreClient()
	if err != nil {
		return err
	}
	defer fs.Close()

	// Snapshot metadata
	snapshotRef := fs.Collection("jira_snapshots").Doc(dateID)
	err = fs.RunTransaction(c.UserContext(), func(ctx context.Context, tx *firestore.Transaction) error {
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

	return c.JSON(fiber.Map{
		"total":   len(issues),
		"success": err == nil,
		"data":    issues,
	})
}

func Me(c *fiber.Ctx) error {
	client, err := njir.NewClient()
	if err != nil {
		return err
	}
	me, _, err := client.User.GetSelf()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"me": me,
	})
}

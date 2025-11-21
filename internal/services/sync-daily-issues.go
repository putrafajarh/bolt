package services

import (
	"context"
	"fmt"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/putrafajarh/bolt/internal/njir"
	firestoreRepo "github.com/putrafajarh/bolt/internal/repository/firestore"
	"github.com/putrafajarh/bolt/internal/repository/postgres"
	"gorm.io/gorm"
)

func (s *JiraServiceImpl) SyncDailyIssues(ctx context.Context, date time.Time) ([]jira.Issue, error) {
	jqlDate := date.Format("2006/01/02")

	sourceStatus := "READY FOR SIT"
	jql := fmt.Sprintf("Status CHANGED TO '%s' ON '%s' ORDER BY assignee asc , created ASC", sourceStatus, jqlDate)

	// Search jira issues that transition to SIT
	issues, err := s.njirClient.SearchIssuesWithTransitions(ctx, jql)
	if err != nil {
		return nil, err
	}

	jiraRepo := firestoreRepo.NewJiraSnapshotRepository(s.firestoreClient)

	// Delete old snapshots older than 30 days
	if err := jiraRepo.CleanupSnapshots(ctx, 30); err != nil {
		return nil, err
	}

	// Store found jira issues to firestore
	_, err = jiraRepo.StoreSnapshots(ctx, date, jql, sourceStatus, issues)
	if err != nil {
		return nil, err
	}

	var users []*jira.User
	var projects []jira.Project

	for _, issue := range issues {
		projects = append(projects, issue.Fields.Project)
		users = append(users, issue.Fields.Assignee)
		users = append(users, issue.Fields.Reporter)
	}

	if err := syncUsers(ctx, users, s.db); err != nil {
		return nil, err
	}

	if err := syncProjects(ctx, projects, s.db); err != nil {
		return nil, err
	}

	if err := syncIssues(ctx, issues, s.db); err != nil {
		return nil, err
	}

	return issues, nil
}

func syncUsers(ctx context.Context, users []*jira.User, db *gorm.DB) error {
	userRepo := postgres.NewJiraUserRepository(db)
	seen := make(map[string]struct{}, len(users))
	for _, user := range users {
		if user != nil {
			key := user.AccountID
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			if err := userRepo.Upsert(ctx, &postgres.User{
				AccountID:   user.AccountID,
				DisplayName: user.DisplayName,
				Email:       user.EmailAddress,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func syncProjects(ctx context.Context, projects []jira.Project, db *gorm.DB) error {
	projectRepo := postgres.NewJiraProjectRepository(db)
	seen := make(map[string]struct{}, len(projects))
	for _, project := range projects {
		if project.ID != "" {
			key := project.ID
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			if err := projectRepo.Upsert(ctx, &postgres.Project{
				ProjectID:   project.ID,
				ProjectName: project.Name,
				ProjectKey:  project.Key,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func syncIssues(ctx context.Context, issues []jira.Issue, db *gorm.DB) error {
	issueRepo := postgres.NewJiraIssueRepository(db)
	for _, issue := range issues {
		if issue.Fields == nil {
			continue
		}

		storyPoint := njir.GetIssueStoryPoint(issue.Fields)
		if issue.ID != "" {

			if err := issueRepo.Upsert(ctx, &postgres.Issue{
				IssueID:    issue.ID,
				IssueKey:   issue.Key,
				Summary:    issue.Fields.Summary,
				StoryPoint: storyPoint,
				ProjectID:  issue.Fields.Project.ID,
				IssueType:  issue.Fields.Type.Name,
				AssigneeID: func() *string {
					if issue.Fields.Assignee == nil {
						return nil
					}
					return &issue.Fields.Assignee.AccountID
				}(),
				ReporterID: func() *string {
					if issue.Fields.Reporter == nil {
						return nil
					}
					return &issue.Fields.Reporter.AccountID
				}(),
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

package postgres

import (
	"context"

	"gorm.io/gorm"
)

type Issue struct {
	gorm.Model
	IssueID    string  `gorm:"column:issue_id;primaryKey" json:"issue_id"`
	IssueKey   string  `gorm:"column:issue_key" json:"issue_key"`
	Summary    string  `gorm:"column:summary" json:"summary"`
	StoryPoint int     `gorm:"column:story_point" json:"story_point"`
	ProjectID  string  `gorm:"column:project_id" json:"project_id"`
	IssueType  string  `gorm:"column:issue_type" json:"issue_type"`
	AssigneeID *string `gorm:"column:assignee_id" json:"assignee_id"`
	Assignee   *User   `gorm:"foreignKey:AssigneeID;references:AccountID"`
	ReporterID *string `gorm:"column:reporter_id" json:"reporter_id"`
	Reporter   *User   `gorm:"foreignKey:ReporterID;references:AccountID"`
}

type JiraIssueRepository struct {
	DB *gorm.DB
}

func NewJiraIssueRepository(db *gorm.DB) *JiraIssueRepository {
	return &JiraIssueRepository{
		DB: db,
	}
}

func (r *JiraIssueRepository) GetOne(ctx context.Context, key string) error {
	return nil
}

func (r *JiraIssueRepository) Upsert(ctx context.Context, issue *Issue) error {
	return r.DB.Save(issue).Error
}

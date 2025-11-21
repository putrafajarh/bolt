package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ProjectID   string `gorm:"column:project_id;primaryKey" json:"project_id"`
	ProjectName string `gorm:"column:project_name" json:"project_name"`
	ProjectKey  string `gorm:"column:project_key" json:"project_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type JiraProjectRepository struct {
	DB *gorm.DB
}

func NewJiraProjectRepository(db *gorm.DB) *JiraProjectRepository {
	return &JiraProjectRepository{
		DB: db,
	}
}

func (r *JiraProjectRepository) Upsert(ctx context.Context, project *Project) error {
	return r.DB.Save(project).Error
}

func (r *JiraProjectRepository) GetByProjectID(ctx context.Context, projectID string) (*Project, error) {
	var project Project
	if err := r.DB.Where("project_id = ?", projectID).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *JiraProjectRepository) All(ctx context.Context) ([]Project, error) {
	var projects []Project
	if err := r.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *JiraProjectRepository) Issues(ctx context.Context, projectID string) ([]Issue, error) {
	var issues []Issue
	if err := r.DB.
		Preload("Assignee").
		Preload("Reporter").
		Where("project_id = ?", projectID).
		Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

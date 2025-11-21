package postgres

import (
	"context"

	"gorm.io/gorm"
)

type JiraSprintRepository struct {
	DB *gorm.DB
}

func NewJiraSprintRepository(db *gorm.DB) *JiraSprintRepository {
	return &JiraSprintRepository{
		DB: db,
	}
}

func (r *JiraSprintRepository) GetOne(ctx context.Context, key string) error {
	return nil
}

package postgres

import (
	"context"

	"gorm.io/gorm"
)

type JiraTransitionRepository struct {
	DB *gorm.DB
}

func NewJiraTransitionRepository(db *gorm.DB) *JiraTransitionRepository {
	return &JiraTransitionRepository{
		DB: db,
	}
}

func (r *JiraTransitionRepository) GetOne(ctx context.Context, key string) error {
	return nil
}

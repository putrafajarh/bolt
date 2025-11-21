package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	AccountID   string `gorm:"column:account_id;primaryKey"`
	DisplayName string `gorm:"column:display_name"`
	Email       string `gorm:"column:email"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type JiraUserRepository struct {
	DB *gorm.DB
}

func NewJiraUserRepository(db *gorm.DB) *JiraUserRepository {
	return &JiraUserRepository{
		DB: db,
	}
}

func (r *JiraUserRepository) GetOne(ctx context.Context, accountID string) (*User, error) {
	var user User
	if err := r.DB.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *JiraUserRepository) Upsert(ctx context.Context, user *User) error {
	return r.DB.Save(user).Error
}

func (r *JiraUserRepository) AssignedIssues(ctx context.Context, accountID string) ([]Issue, error) {
	var issues []Issue
	if err := r.DB.Where("assignee_id = ?", accountID).Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

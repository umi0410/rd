package repository

import (
	"context"

	"gorm.io/gorm"
	"rd/entity"
)

func NewGormAuthRepository(db *gorm.DB) *GormAuthRepository {
	repository := &GormAuthRepository{cli: db}

	return repository
}

type GormAuthRepository struct {
	cli *gorm.DB
}

func (r *GormAuthRepository) GetUser(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	if err := r.cli.Model(&entity.User{}).
		Preload("Groups").
		First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

package service

import (
	"context"

	"rd/repository"
	"rd/util"
)

type AuthService interface {
	IsAdmin(ctx context.Context, username string) bool
	IsInGroup(ctx context.Context, username, group string) bool
}

type AuthServiceImpl struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{repo: repo}
}

func (a *AuthServiceImpl) IsAdmin(ctx context.Context, username string) bool {
	return a.IsInGroup(ctx, username, "admin")
}

func (a *AuthServiceImpl) IsInGroup(ctx context.Context, username, group string) bool {
	log := util.GetLogger(ctx)
	user, err := a.repo.GetUser(ctx, username)
	if err != nil {
		log.Errorf("There was an error retrieving the user(%s), so just judging the user is not the group(%s): %+v ", username, group, err)
		return false
	}

	for _, g := range user.Groups {
		if g.Name == group {
			return true
		}
	}
	return false
}

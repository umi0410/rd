package service

import (
	"fmt"

	"github.com/pkg/errors"
	"rd/domain"
	"rd/entity"
	"rd/mapper"
	"rd/repository"
)

type AliasService interface {
	Create(*entity.Alias) (*domain.Alias, error)
	List() ([]*domain.Alias, error)
	ListByGroup(group string) ([]*domain.Alias, error)
	ListByGroupAndAlias(group, alias string) ([]*domain.Alias, error)
	Update(*entity.Alias) (*domain.Alias, error)
	Delete(id int) (*domain.Alias, error)
}

type AliasServiceImpl struct {
	repo repository.AliasRepository
}

var (
	ErrDuplicatedAlias = fmt.Errorf("duplicated aliases already exist")
)

func (s AliasServiceImpl) Create(alias *entity.Alias) (*domain.Alias, error) {
	if len(s.repo.ListByGroupAndAlias(alias.AliasGroup, alias.Name)) != 0 {
		return nil, errors.WithStack(ErrDuplicatedAlias)
	}

	alias, err := s.repo.Create(alias)
	if err != nil {
		return nil, err
	}

	return mapper.AliasFromEntityToDomain(alias), nil
}

func (s AliasServiceImpl) List() ([]*domain.Alias, error) {
	aliases := s.repo.List()

	return mapper.AliasesFromEntityToDomain(aliases), nil
}

func (s AliasServiceImpl) ListByGroup(group string) ([]*domain.Alias, error) {
	aliases := s.repo.ListByGroup(group)

	return mapper.AliasesFromEntityToDomain(aliases), nil
}

func (s AliasServiceImpl) ListByGroupAndAlias(group, alias string) ([]*domain.Alias, error) {
	aliases := s.repo.ListByGroupAndAlias(group, alias)

	return mapper.AliasesFromEntityToDomain(aliases), nil
}

func (s AliasServiceImpl) Update(*entity.Alias) (*domain.Alias, error) {
	panic("implement me")
}

func (s AliasServiceImpl) Delete(id int) (*domain.Alias, error) {
	alias, err := s.repo.Delete(id)
	if err != nil {
		return nil, err
	}

	return mapper.AliasFromEntityToDomain(alias), nil
}

func NewAliasService(repo repository.AliasRepository) AliasService {
	return &AliasServiceImpl{
		repo: repo,
	}
}

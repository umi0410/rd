package repository

import (
	"rd/domain"
)

type LocalAliasRepository struct {
	Store Store
}

func (repo *LocalAliasRepository) List() []*domain.Alias {
	return repo.Store.List()
}

func (repo *LocalAliasRepository) ListByAlias(alias string) []*domain.Alias {
	return repo.Store.ListByAlias(alias)
}

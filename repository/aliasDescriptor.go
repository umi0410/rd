package repository

import (
	"rd/domain"
)

type LocalAliasDescriptorRepository struct {
	Store Store
}

func (repo *LocalAliasDescriptorRepository) List() []*domain.AliasDescriptor {
	return repo.Store.List()
}

func (repo *LocalAliasDescriptorRepository) ListByAlias(alias string) []*domain.AliasDescriptor {
	return repo.Store.ListByAlias(alias)
}

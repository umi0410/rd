package repository

import "mytot/domain"

type Store interface {
	Add(*domain.AliasDescriptor) error
	Delete() error
	List() []*domain.AliasDescriptor
	ListByAlias(alias string) []*domain.AliasDescriptor
	Reload() error
}

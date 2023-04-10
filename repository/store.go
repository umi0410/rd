package repository

import "rd/domain"

type Store interface {
	Add(*domain.Alias) error
	Delete() error
	List() []*domain.Alias
	ListByAlias(alias string) []*domain.Alias
	Reload() error
}

package service

import "rd/domain"

type AliasRepository interface {
	Create(*domain.Alias) (*domain.Alias, error)
	List() []*domain.Alias
	ListByGroup(group string) []*domain.Alias
	ListByGroupAndAlias(group, alias string) []*domain.Alias
}

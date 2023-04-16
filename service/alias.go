package service

import "rd/domain"

type AliasRepository interface {
	List() []*domain.Alias
	ListByGroup(group string) []*domain.Alias
	ListByGroupAndAlias(group, alias string) []*domain.Alias
}

package service

import "rd/domain"

type AliasRepository interface {
	List() []*domain.Alias
	ListByAlias(alias string) []*domain.Alias
}

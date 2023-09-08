package repository

import (
	"rd/entity"
)

type AliasRepository interface {
	Create(*entity.Alias) (*entity.Alias, error)
	List() []*entity.Alias
	ListByGroup(group string) []*entity.Alias
	ListByGroupAndAlias(group, alias string) []*entity.Alias
	Close() error
}

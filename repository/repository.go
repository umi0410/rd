package repository

import (
	"time"

	"rd/entity"
)

type AliasRepository interface {
	Create(*entity.Alias) (*entity.Alias, error)
	List() []*entity.Alias
	ListByGroup(group string) []*entity.Alias
	ListByGroupAndAlias(group, alias string) []*entity.Alias
	Delete(id int) (*entity.Alias, error)
	Close() error
}

type EventAliasHitRepository interface {
	Create(evt *entity.EventAliasHit) (*entity.EventAliasHit, error)
	ListByAliasIds(aliasIds []uint) []*entity.EventAliasHit
	ListByAliasIdsAndGreaterThanCreatedAt(aliasIds []uint, createdAt time.Time) []*entity.EventAliasHit
}

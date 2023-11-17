package repository

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rd/entity"
)

type GormEventAliasHitRepository struct {
	cli *gorm.DB
}

func NewGormEventAliasHitRepository(db *gorm.DB) *GormEventAliasHitRepository {
	return &GormEventAliasHitRepository{cli: db}
}

func (r *GormEventAliasHitRepository) Create(evt *entity.EventAliasHit) (*entity.EventAliasHit, error) {
	res := r.cli.Create(evt)
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return evt, nil
}

func (r *GormEventAliasHitRepository) ListByAliasIds(aliasIds []uint) []*entity.EventAliasHit {
	events := make([]*entity.EventAliasHit, 0, 32)
	res := r.cli.Find(&events, aliasIds)
	if res.Error != nil {
		log.Errorf("%+v", errors.WithStack(res.Error))
		return []*entity.EventAliasHit{}
	}

	return events
}

func (r *GormEventAliasHitRepository) ListByAliasIdsAndGreaterThanCreatedAt(aliasIds []uint, createdAt time.Time) []*entity.EventAliasHit {
	events := make([]*entity.EventAliasHit, 0, 32)
	res := r.cli.Where("created_at >= ? AND alias_fk IN ?", createdAt, aliasIds).Find(&events)
	if res.Error != nil {
		log.Errorf("%+v", errors.WithStack(res.Error))
		return []*entity.EventAliasHit{}
	}

	return events
}

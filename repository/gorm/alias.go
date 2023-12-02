package gorm

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"rd/config"

	"gorm.io/gorm"
	"rd/entity"
)

type GormAliasRepository struct {
	cli *gorm.DB
}

func NewGormAliasRepository(db *gorm.DB, kind config.RepositoryKind) *GormAliasRepository {
	repository := &GormAliasRepository{cli: db}

	// XXX: when using sqlite memory, automatically
	// add fixture data
	if kind == config.RepoKindSqliteMemory {
		repository.fixture()
	}

	return repository
}

func (r *GormAliasRepository) Create(alias *entity.Alias) (*entity.Alias, error) {
	res := r.cli.Create(alias)
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return alias, nil
}

func (r *GormAliasRepository) List() []*entity.Alias {
	aliases := make([]*entity.Alias, 0, 32)
	res := r.cli.Find(&aliases)
	if res.Error != nil {
		log.Errorf("%+v", errors.WithStack(res.Error))
		return aliases
	}

	return aliases
}

func (r *GormAliasRepository) ListByGroup(group string, recentHitCountSince time.Time) []*entity.Alias {
	aliases := make([]*entity.Alias, 0, 32)
	res := r.cli.Select("aliases.*, COALESCE(hit_count.count, 0)").
		Where("alias_group = ?", group).
		Joins("LEFT JOIN (SELECT alias_fk, count(*) AS count FROM event_alias_hits WHERE ? < event_alias_hits.created_at GROUP BY alias_fk) AS hit_count ON id = alias_fk", recentHitCountSince).
		Order("hit_count.count desc").Find(&aliases)
	if res.Error != nil {
		log.Errorf("%+v", errors.WithStack(res.Error))
		return aliases
	}

	return aliases
}

func (r *GormAliasRepository) ListByGroupAndAlias(group, alias string) []*entity.Alias {
	aliases := make([]*entity.Alias, 0, 32)
	res := r.cli.Model(new(entity.Alias)).
		Where("alias_group = ? AND name = ?", group, alias).
		Find(&aliases)
	if res.Error != nil {
		log.Errorf("%+v", errors.WithStack(res.Error))
		return aliases
	}

	return aliases
}

func (r *GormAliasRepository) Get(id int) (*entity.Alias, error) {
	alias := &entity.Alias{}
	res := r.cli.First(alias, id)
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return alias, nil
}

func (r *GormAliasRepository) Delete(id int) (*entity.Alias, error) {
	alias := &entity.Alias{}
	res := r.cli.First(alias, id)
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	res = r.cli.Delete(&entity.Alias{}, id)
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return alias, nil
}

func (r *GormAliasRepository) Close() error {
	db, err := r.cli.DB()
	if err != nil {
		return errors.WithStack(err)
	}
	err = db.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Infof("Closed DB")

	return nil
}

func (r *GormAliasRepository) fixture() {
	for _, f := range []*entity.Alias{{
		AliasGroup:  "james",
		Name:        "naver",
		Destination: "https://naver.com",
	}, {
		AliasGroup:  "james",
		Name:        "google",
		Destination: "https://google.com",
	}, {
		AliasGroup:  "james",
		Name:        "github",
		Destination: "https://github.com",
	},
	} {
		_, err := r.Create(f)
		if err != nil {
			log.Panicf("%+v", err)
		}
	}
}

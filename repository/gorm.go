package repository

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"rd/config"

	//"gorm.io/driver/mysql"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"rd/entity"
)

type GormAliasRepository struct {
	cli *gorm.DB
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

func (r *GormAliasRepository) ListByGroup(group string) []*entity.Alias {
	aliases := make([]*entity.Alias, 0, 32)
	res := r.cli.Where("alias_group = ?", group).Find(&aliases)
	if res.Error != nil {
		log.Errorf("%+v", errors.WithStack(res.Error))
		return aliases
	}

	return aliases
}

func (r *GormAliasRepository) ListByGroupAndAlias(group, alias string) []*entity.Alias {
	aliases := make([]*entity.Alias, 0, 32)
	res := r.cli.Where("alias_group = ? AND name = ?", group, alias).Find(&aliases)
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

//
//func (r *GormAliasRepository) ListByAlias(alias string) []*domain.Alias {
//	aliases := make([]*domain.Alias, 0, 32)
//	res := r.cli.Where("name = ?", alias).Find(aliases)
//	if res.Error != nil {
//		log.Errorf("%+v", errors.WithStack(res.Error))
//		return aliases
//	}
//
//	return aliases
//}
//
//func (*GormAliasRepository) Reload() error {
//	//TODO implement me
//	panic("implement me")
//}

// dns: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func NewGormRepository(kind config.RepositoryKind, dsn string) (AliasRepository, EventAliasHitRepository, error) {
	var (
		db  *gorm.DB
		err error
	)
	switch kind {
	case config.RepoKindMysql:
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
	case config.RepoKindSqlite:
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
	}

	if err := db.AutoMigrate(&entity.Alias{}, &entity.EventAliasHit{}); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	aliasRepo, err := &GormAliasRepository{
		cli: db,
	}, nil
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	eventAliasRepo, err := &GormEventAliasHitRepository{
		cli: db,
	}, nil
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	// XXX: when using sqlite memory, automatically
	// add fixture data
	if kind == config.RepoKindSqliteMemory {
		aliasRepo.fixture()
	}

	return aliasRepo, eventAliasRepo, nil
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

type GormEventAliasHitRepository struct {
	cli *gorm.DB
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

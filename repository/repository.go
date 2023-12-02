package repository

import (
	"context"
	"time"

	sqliteDriver "github.com/glebarez/sqlite"
	"github.com/pkg/errors"
	mysqlDriver "gorm.io/driver/mysql"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"rd/config"
	"rd/entity"
)

type AliasRepository interface {
	Create(*entity.Alias) (*entity.Alias, error)
	List() []*entity.Alias
	ListByGroup(group string, recentHitCountSince time.Time) []*entity.Alias
	ListByGroupAndAlias(group, alias string) []*entity.Alias
	Delete(id int) (*entity.Alias, error)
	Close() error
}

type EventAliasHitRepository interface {
	Create(evt *entity.EventAliasHit) (*entity.EventAliasHit, error)
	ListByAliasIds(aliasIds []uint) []*entity.EventAliasHit
	ListByAliasIdsAndGreaterThanCreatedAt(aliasIds []uint, createdAt time.Time) []*entity.EventAliasHit
}

type AuthRepository interface {
	GetUser(ctx context.Context, username string) (*entity.User, error)
}

// dns: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func NewDB(kind config.RepositoryKind, dsn string) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	switch kind {
	case config.RepoKindMysql:
		db, err = gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	case config.RepoKindCockroachdb:
		db, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	case config.RepoKindSqlite:
		logger := gormLogger.Default
		logger.LogMode(gormLogger.Info)
		db, err = gorm.Open(sqliteDriver.Open(dsn), &gorm.Config{Logger: logger})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// XXX: Uncomment these lines only if you want to enable auto-migrate
	//if err := db.AutoMigrate(&entity.Alias{}, &entity.EventAliasHit{}, &entity.Group{}, &entity.User{}); err != nil {
	//	return nil, errors.WithStack(err)
	//}
	return db, nil
}

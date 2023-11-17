package cmd

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rd/config"
	"rd/repository"
	"rd/service"
)

func initialize() (repository.AliasRepository, repository.EventAliasHitRepository, service.AliasService, service.AuthService) {
	repoCfg := config.Cfg.Repository
	var (
		db                *gorm.DB
		authRepo          repository.AuthRepository
		aliasRepo         repository.AliasRepository
		eventAliasHitRepo repository.EventAliasHitRepository
		authService       service.AuthService
		aliasService      service.AliasService
		err               error
	)
	switch repoCfg.Kind {
	case config.RepoKindSqliteMemory, config.RepoKindMysql, config.RepoKindCockroachdb:
		dsn := repoCfg.SqliteMemory.Dsn
		if repoCfg.Kind == config.RepoKindMysql {
			dsn = repoCfg.Mysql.Dsn
		} else if repoCfg.Kind == config.RepoKindCockroachdb {
			dsn = repoCfg.Cockroachdb.Dsn
		}
		if db, err = repository.NewDB(repoCfg.Kind, dsn); err != nil {
			log.Panicf("%+v", err)
		}
		aliasRepo = repository.NewGormAliasRepository(db, repoCfg.Kind)
		eventAliasHitRepo = repository.NewGormEventAliasHitRepository(db)
		authRepo = repository.NewGormAuthRepository(db)
	default:
		log.Panicf("Unsupported repo kind")
	}

	authService = service.NewAuthService(authRepo)
	aliasService = service.NewAliasService(aliasRepo, eventAliasHitRepo, authService)

	return aliasRepo, eventAliasHitRepo, aliasService, authService
}

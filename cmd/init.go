package cmd

import (
	"context"

	"github.com/pkg/errors"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
	gormgorm "gorm.io/gorm"
	"rd/config"
	"rd/repository"
	"rd/repository/dynamodb"
	"rd/repository/gorm"
	"rd/service"
)

func initialize() (repository.AliasRepository, repository.EventAliasHitRepository, service.AliasService, service.AuthService) {
	repoCfg := config.Cfg.Repository
	var (
		db                *gormgorm.DB
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
		aliasRepo = gorm.NewGormAliasRepository(db, repoCfg.Kind)
		eventAliasHitRepo = gorm.NewGormEventAliasHitRepository(db)
		authRepo = gorm.NewGormAuthRepository(db)
	case config.RepoKindDynamodb:
		awsConfig, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(repoCfg.Dynamodb.Region))
		if err != nil {
			log.Panicf("%+v", errors.WithStack(err))
		}
		aliasRepo = dynamodb.NewDynamodbAliasRepository(repoCfg.Dynamodb, awsConfig)
		if err != nil {
			log.Panicf("%+v", err)
		}
		authRepo = dynamodb.NewDynamodbAuthRepository(repoCfg.Dynamodb, awsConfig)
	default:
		log.Panicf("Unsupported repo kind")
	}

	authService = service.NewAuthService(authRepo)
	aliasService = service.NewAliasService(aliasRepo, eventAliasHitRepo, authService)

	return aliasRepo, eventAliasHitRepo, aliasService, authService
}

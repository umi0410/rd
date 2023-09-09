package cmd

import (
	log "github.com/sirupsen/logrus"
	"rd/config"
	"rd/repository"
	"rd/service"
)

func initialize() (repository.AliasRepository, repository.EventAliasHitRepository, service.AliasService) {
	repoCfg := config.Cfg.Repository
	var (
		aliasRepo         repository.AliasRepository
		eventAliasHitRepo repository.EventAliasHitRepository
		aliasSvc          service.AliasService
		err               error
	)
	switch repoCfg.Kind {
	case config.RepoKindSqliteMemory, config.RepoKindMysql:
		dsn := repoCfg.SqliteMemory.Dsn
		if repoCfg.Kind == config.RepoKindMysql {
			dsn = repoCfg.Mysql.Dsn
		}
		aliasRepo, eventAliasHitRepo, err = repository.NewGormRepository(repoCfg.Kind, dsn)
		if err != nil {
			log.Panicf("%+v", err)
		}
	//case config.RepositoryKindNats:
	//	repo, err = repository.NewNatsRepository(repository.NatsRepositoryConfig{
	//		Host:     repoCfg.Nats.Host,
	//		Port:     repoCfg.Nats.Port,
	//		Username: repoCfg.Nats.Username,
	//		Password: repoCfg.Nats.Password,
	//		Bucket:   repoCfg.Nats.Bucket,
	//	})
	//	if err != nil {
	//		log.Panicf("%+v", err)
	//	}
	//}
	default:
		log.Panicf("Unsupported repo kind")
	}

	aliasSvc = service.NewAliasService(aliasRepo, eventAliasHitRepo)

	return aliasRepo, eventAliasHitRepo, aliasSvc
}

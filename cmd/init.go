package cmd

import (
	log "github.com/sirupsen/logrus"
	"rd/config"
	"rd/repository"
)

func initialize() repository.AliasRepository {
	repoCfg := config.Cfg.Repository
	var repo repository.AliasRepository
	var err error
	switch repoCfg.Kind {
	case config.RepoKindSqliteMemory, config.RepoKindMysql:
		dsn := repoCfg.SqliteMemory.Dsn
		if repoCfg.Kind == config.RepoKindMysql {
			dsn = repoCfg.Mysql.Dsn
		}
		repo, err = repository.NewGormRepository(repoCfg.Kind, dsn)
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

	return repo
}

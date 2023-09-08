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
	case config.RepoKindSqliteMemory:
		repo, err = repository.NewSqlLiteRepository(repoCfg.SqliteMemory.Dsn)
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

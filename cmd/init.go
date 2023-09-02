package cmd

import (
	log "github.com/sirupsen/logrus"
	"rd/config"
	"rd/repository"
	"rd/service"
)

func initialize() service.AliasRepository {
	repoCfg := config.Cfg.Repository
	repo, err := repository.NewNatsRepository(repository.NatsRepositoryConfig{
		Host:     repoCfg.Nats.Host,
		Port:     repoCfg.Nats.Port,
		Username: repoCfg.Nats.Username,
		Password: repoCfg.Nats.Password,
		Bucket:   repoCfg.Nats.Bucket,
	})
	if err != nil {
		log.Panicf("%+v", err)
	}

	return repo
}

package config

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	Repository struct {
		Kind         RepositoryKind
		SqliteMemory struct {
			Dsn string
		}
		Nats struct {
			Host     string
			Port     int
			Username string
			Password string
			Bucket   string
		}
	}
}

type RepositoryKind string

const (
	RepoKindSqlite       RepositoryKind = "sqlite"
	RepoKindSqliteMemory RepositoryKind = "sqliteMemory"
	RepositoryKindNats   RepositoryKind = "nats"
)

func Load() error {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.AddConfigPath(os.ExpandEnv("${HOME}/.config/rd"))
	vp.AddConfigPath(os.ExpandEnv("./config"))

	// configName is used before viper so this parameter
	// is parsed by not viper but os.Getenv.
	configName := os.Getenv("RD_CONFIG_NAME")
	if len(configName) == 0 {
		configName = "default"
	}
	vp.SetConfigName(configName)

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// skip
			log.Warnf("No viper config file")
		} else {
			return errors.WithStack(err)
		}
	} else {
		log.Infof("Found a viper config file")
	}

	Cfg = new(Config)
	if err := vp.Unmarshal(Cfg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

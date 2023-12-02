package config

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Cfg *Config

type TlsConfig struct {
	Mode             TlsMode
	CertFile         string
	PrivateKeyFile   string
	ClientCaCertFile string
}

type DynamodbConfig struct {
	Region                   string
	AliasTableName           string
	AliasTablePkName         string
	UserTableName            string
	UserTablePkName          string
	AliasHitEventTableName   string
	AliasHitEventTablePkName string
	//Profile string
}

type Config struct {
	Tls        TlsConfig
	Repository struct {
		Kind         RepositoryKind
		SqliteMemory struct {
			Dsn string
		}
		Mysql struct {
			Dsn string
		}
		Cockroachdb struct {
			Dsn string
		}
		Nats struct {
			Host     string
			Port     int
			Username string
			Password string
			Bucket   string
		}
		Dynamodb DynamodbConfig
	}
}

type RepositoryKind string

const (
	RepoKindSqlite       RepositoryKind = "sqlite"
	RepoKindSqliteMemory RepositoryKind = "sqliteMemory"
	RepoKindMysql        RepositoryKind = "mysql"
	RepoKindCockroachdb  RepositoryKind = "cockroachdb"
	RepoKindDynamodb     RepositoryKind = "dynamodb"
	RepositoryKindNats   RepositoryKind = "nats"
)

type TlsMode string

const (
	TlsModePlaintext TlsMode = "plaintext"
	TlsModeSimple    TlsMode = "simple"
	TlsModeMutual    TlsMode = "mutual"
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

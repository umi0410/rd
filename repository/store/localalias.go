package store

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"rd/domain"
)

type LocalStore struct {
	vp      *viper.Viper
	config  *config
	aliases map[string][]*domain.Alias
}

type config struct {
	Aliases []*domain.Alias
}

func NewLocalStore() (*LocalStore, error) {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.AddConfigPath(os.ExpandEnv("${HOME}/.config/rd"))
	vp.AddConfigPath(os.ExpandEnv("./config"))
	configName := os.Getenv("RD_CONFIG_NAME")
	if configName == "" {
		configName = "default"
	}
	vp.SetConfigName(configName)

	store := &LocalStore{
		config:  new(config),
		vp:      vp,
		aliases: map[string][]*domain.Alias{},
	}

	if err := store.Reload(); err != nil {
		return nil, errors.WithStack(err)
	}

	return store, nil
}

func (s *LocalStore) Add(aliasDescriptor *domain.Alias) error {
	s.config.Aliases = append(s.config.Aliases, aliasDescriptor)
	viper.Set("config.aliases", s.config.Aliases)
	if err := viper.WriteConfig(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *LocalStore) List() []*domain.Alias {
	var ret []*domain.Alias
	for _, dscList := range s.aliases {
		for _, dsc := range dscList {
			ret = append(ret, dsc)
		}
	}

	if ret == nil {
		ret = []*domain.Alias{}
	}

	return ret
}

func (s *LocalStore) ListByAlias(alias string) []*domain.Alias {
	ret := s.aliases[alias]
	if ret == nil {
		ret = []*domain.Alias{}
	}

	return ret
}

func (s *LocalStore) Delete() error {
	//TODO implement me
	panic("implement me")
}

func (s *LocalStore) Reload() error {
	if err := s.vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// skip
			log.Warnf("No viper config file")
		} else {
			return errors.WithStack(err)
		}
	} else {
		log.Infof("Found a viper config file")
	}

	if err := s.vp.Unmarshal(s.config); err != nil {
		return errors.WithStack(err)
	}

	for _, dsc := range s.config.Aliases {
		s.aliases[dsc.Name] = append(s.aliases[dsc.Name], dsc)
	}

	return nil
}

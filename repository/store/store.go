package store

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"rd/domain"
)

type LocalStore struct {
	vp               *viper.Viper
	config           *config
	aliasDescriptors map[string][]*domain.AliasDescriptor
}

type config struct {
	AliasDescriptors []*domain.AliasDescriptor
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
		config:           new(config),
		vp:               vp,
		aliasDescriptors: map[string][]*domain.AliasDescriptor{},
	}

	if err := store.Reload(); err != nil {
		return nil, errors.WithStack(err)
	}

	return store, nil
}

func (s *LocalStore) Add(aliasDescriptor *domain.AliasDescriptor) error {
	s.config.AliasDescriptors = append(s.config.AliasDescriptors, aliasDescriptor)
	viper.Set("config.aliasDescriptors", s.config.AliasDescriptors)
	if err := viper.WriteConfig(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *LocalStore) List() []*domain.AliasDescriptor {
	var ret []*domain.AliasDescriptor
	for _, dscList := range s.aliasDescriptors {
		for _, dsc := range dscList {
			ret = append(ret, dsc)
		}
	}

	if ret == nil {
		ret = []*domain.AliasDescriptor{}
	}

	return ret
}

func (s *LocalStore) ListByAlias(alias string) []*domain.AliasDescriptor {
	ret := s.aliasDescriptors[alias]
	if ret == nil {
		ret = []*domain.AliasDescriptor{}
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

	for _, dsc := range s.config.AliasDescriptors {
		s.aliasDescriptors[dsc.Alias] = append(s.aliasDescriptors[dsc.Alias], dsc)
	}

	return nil
}

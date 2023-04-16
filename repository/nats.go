package repository

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"rd/domain"
	"strconv"
	"strings"
)

type NatsRepository struct {
	conn       *nats.Conn
	aliases    map[string][]*domain.Alias
	objStore   nats.ObjectStore
	objWatcher nats.ObjectWatcher
}

type NatsRepositoryConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Bucket   string
}

func NewNatsRepository(config NatsRepositoryConfig) (*NatsRepository, error) {
	aliases := map[string][]*domain.Alias{}
	conn, err := nats.Connect(fmt.Sprintf("nats://%s:%s", config.Host, strconv.Itoa(config.Port)),
		nats.UserInfo(config.Username, config.Password))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jet, err := conn.JetStream()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	objStore, err := jet.ObjectStore(config.Bucket)
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			log.Warnf("Try creating a new bucket name of which is %s", config.Bucket)
		}
		return nil, errors.WithStack(err)
	}

	objWatcher, err := objStore.Watch()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	repo := &NatsRepository{
		conn:       conn,
		aliases:    aliases,
		objStore:   objStore,
		objWatcher: objWatcher,
	}

	go func() {
		err := repo.watch()
		log.Panicf("%+v", errors.WithStack(err))
	}()

	return repo, nil
}

func (repo *NatsRepository) watch() error {
	var err error
	defer func() {
		err = errors.WithStack(repo.objWatcher.Stop())
		log.Errorf("%+v", err)
	}()

	for objInfo := range repo.objWatcher.Updates() {
		// objInfo.Name means the name of the file in the bucket which is updated.
		if objInfo == nil {
			log.Infof("pass")
			continue
		}
		log.Infof("Object watcher got updated(Digest=%s, Bucket=%s, Object=%s, Deleted=%t), ", objInfo.Digest, objInfo.Bucket, objInfo.Name, objInfo.Deleted)
		// cut off `.yaml`.
		group := strings.TrimSuffix(objInfo.Name, ".yaml")

		if objInfo.Deleted {
			delete(repo.aliases, group)
			continue
		}
		// name argument means the name of the file in the bucket
		// file argument means the name of the local file as the destination
		//repo.objStore.GetFile(objInfo.Name, objInfo.Name)
		res, err := repo.objStore.Get(objInfo.Name)

		if err != nil {
			log.Errorf("%+v", errors.WithStack(err))
			continue
		}

		data, err := io.ReadAll(res)
		if err != nil {
			log.Errorf("%+v", errors.WithStack(err))
			continue
		}

		aliases := make([]*domain.Alias, 0)
		err = yaml.Unmarshal(data, &aliases)
		if err != nil {
			log.Errorf("%+v", errors.WithStack(err))
			continue
		}

		for _, a := range aliases {
			a.Group = group
		}

		repo.aliases[group] = aliases

		// TODO: Log the comparison between before and after
		log.Infof("The update(Digest=%s, Bucket=%s, Object=%s, Deleted=%t) has been applied to the repository.", objInfo.Digest, objInfo.Bucket, objInfo.Name, objInfo.Deleted)
	}

	return err
}

// List all aliases without filtering by Group
func (repo *NatsRepository) List() []*domain.Alias {
	ret := make([]*domain.Alias, 0)
	for _, aliases := range repo.aliases {
		ret = append(ret, aliases...)
	}

	return ret
}

// ListByGroup aliases which are only correspondent to the specific Group
func (repo *NatsRepository) ListByGroup(group string) []*domain.Alias {
	ret, ok := repo.aliases[group]
	if !ok {
		log.Errorf("No such group yet: %s", group)
		return []*domain.Alias{}
	}

	return ret
}

// ListByGroupAndAlias aliases which are only correspondent to the specific Group and Alias
func (repo *NatsRepository) ListByGroupAndAlias(group, alias string) []*domain.Alias {
	aliases, ok := repo.aliases[group]
	if !ok {
		log.Errorf("No such group yet: %s", group)
		return []*domain.Alias{}
	}

	var ret []*domain.Alias
	for _, a := range aliases {
		if a.Name == alias {
			ret = append(ret, a)
		}
	}

	return ret
}

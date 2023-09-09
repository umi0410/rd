package service

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"rd/domain"
	"rd/entity"
	"rd/mapper"
	"rd/repository"
)

type AliasService interface {
	Create(*entity.Alias) (*domain.Alias, error)
	List() ([]*domain.Alias, error)
	// 최근 X time.Duration 동안의 hit count table을 sum한 값을 내림차순으로 정렬한
	// alias row들을 N개 조회하라.
	// SELECT
	ListByGroup(group string) ([]*domain.Alias, error)
	ListByGroupAndAlias(group, alias string) ([]*domain.Alias, error)
	GoTo(group, alias string) ([]*domain.Alias, error)
	Update(*entity.Alias) (*domain.Alias, error)
	Delete(id int) (*domain.Alias, error)
}

type AliasServiceImpl struct {
	repo              repository.AliasRepository
	eventAliasHitRepo repository.EventAliasHitRepository
}

var (
	ErrDuplicatedAlias = fmt.Errorf("duplicated aliases already exist")
)

func (s AliasServiceImpl) Create(alias *entity.Alias) (*domain.Alias, error) {
	if len(s.repo.ListByGroupAndAlias(alias.AliasGroup, alias.Name)) != 0 {
		return nil, errors.WithStack(ErrDuplicatedAlias)
	}

	alias, err := s.repo.Create(alias)
	if err != nil {
		return nil, err
	}

	return mapper.AliasFromEntityToDomain(alias), nil
}

func (s AliasServiceImpl) List() ([]*domain.Alias, error) {
	aliasEntities := s.repo.List()
	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (s AliasServiceImpl) ListByGroup(group string) ([]*domain.Alias, error) {
	if group == "" {
		return s.List()
	}

	aliasEntities := s.repo.ListByGroup(group)
	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (s AliasServiceImpl) ListByGroupAndAlias(group, alias string) ([]*domain.Alias, error) {
	if group == "" && alias == "" {
		return s.List()
	}

	aliasEntities := s.repo.ListByGroupAndAlias(group, alias)
	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (s AliasServiceImpl) GoTo(group, alias string) ([]*domain.Alias, error) {
	aliasEntities := s.repo.ListByGroupAndAlias(group, alias)

	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}
	hit := false
	if len(aliases) >= 1 {
		hit = true
	}
	for _, a := range aliases {
		evt, err := s.eventAliasHitRepo.Create(&entity.EventAliasHit{
			Hit:     hit,
			AliasFK: a.ID,
			User:    "",
		})
		if err != nil {
			log.Errorf("%+v", errors.WithStack(err))
		} else {
			log.Infof("Created an EventAliasHit(%+v)", evt)
		}
	}

	return aliases, nil
}

func (s AliasServiceImpl) Update(*entity.Alias) (*domain.Alias, error) {
	panic("implement me")
}

func (s AliasServiceImpl) Delete(id int) (*domain.Alias, error) {
	aliasEntity, err := s.repo.Delete(id)
	if err != nil {
		return nil, err
	}

	alias := mapper.AliasFromEntityToDomain(aliasEntity)
	if err := s.setRecentHits([]*domain.Alias{alias}); err != nil {
		return nil, err
	}

	return alias, nil
}

func NewAliasService(repo repository.AliasRepository, eventAliasHitRepo repository.EventAliasHitRepository) AliasService {
	return &AliasServiceImpl{
		repo:              repo,
		eventAliasHitRepo: eventAliasHitRepo,
	}
}

func (s AliasServiceImpl) setRecentHits(aliases []*domain.Alias) error {
	var (
		//aliasIds       []uint
		recentDuration     = time.Hour * 24 * 14 // 2 Weeks
		recentDurationName = "Recent 2 Weeks"
	)
	after := time.Now().Add(-recentDuration)
	// TODO: Too many DB requests. Introduce batch queries.
	for _, a := range aliases {
		//aliasIds = append(aliasIds, a.ID)
		events := s.eventAliasHitRepo.ListByAliasIdsAndGreaterThanCreatedAt([]uint{a.ID}, after)
		a.RecentHitCounts[recentDurationName] = len(events)
	}

	return nil
}

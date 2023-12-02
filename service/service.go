package service

import (
	"context"
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
	Create(ctx context.Context, alias *entity.Alias) (*domain.Alias, error)
	List(ctx context.Context, user string) ([]*domain.Alias, error)
	// 최근 X time.Duration 동안의 hit count table을 sum한 값을 내림차순으로 정렬한
	// alias row들을 N개 조회하라.
	// SELECT
	ListByGroup(ctx context.Context, user, group string, sort Sort) ([]*domain.Alias, error)
	ListByGroupAndAlias(ctx context.Context, user, group, alias string, sort Sort) ([]*domain.Alias, error)
	GoTo(ctx context.Context, user, group, alias string) ([]*domain.Alias, error)
	Update(ctx context.Context, alias *entity.Alias) (*domain.Alias, error)
	Delete(ctx context.Context, id int) (*domain.Alias, error)
}

type AliasServiceImpl struct {
	authService       AuthService
	repo              repository.AliasRepository
	eventAliasHitRepo repository.EventAliasHitRepository
}

type Sort string

var (
	SortByDefault         Sort = ""
	SortByRecentHitCounts Sort = "recent_hit_counts_desc"

	ErrNoPermission = fmt.Errorf("no permission")
)

func (s *Sort) Validate() error {
	for _, sort := range []Sort{SortByDefault, SortByRecentHitCounts} {
		if *s == sort {
			return nil
		}
	}
	return errors.New("unknown sort")
}

var (
	ErrDuplicatedAlias = fmt.Errorf("duplicated aliases already exist")
)

func (s AliasServiceImpl) Create(ctx context.Context, alias *entity.Alias) (*domain.Alias, error) {
	if err := alias.Validate(); err != nil {
		return nil, err
	}
	alias.CreatedAt = time.Time{}
	now := time.Now()
	SetWhenZeroTimeValue(&alias.CreatedAt, now)
	SetWhenZeroTimeValue(&alias.UpdatedAt, now)

	if len(s.repo.ListByGroupAndAlias(alias.AliasGroup, alias.Name)) != 0 {
		return nil, errors.WithStack(ErrDuplicatedAlias)
	}

	alias, err := s.repo.Create(alias)
	if err != nil {
		return nil, err
	}

	return mapper.AliasFromEntityToDomain(alias), nil
}

func (s AliasServiceImpl) List(ctx context.Context, user string) ([]*domain.Alias, error) {
	if !s.authService.IsAdmin(ctx, user) {
		return nil, errors.Wrap(ErrNoPermission, "Only admins can list aliases without a group filter.")
	}
	aliasEntities := s.repo.List()
	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (s AliasServiceImpl) ListByGroup(ctx context.Context, user, group string, sort Sort) ([]*domain.Alias, error) {
	if !s.authService.IsInGroup(ctx, user, group) && !s.authService.IsAdmin(ctx, user) {
		return nil, errors.Wrap(ErrNoPermission, fmt.Sprintf("user(%s) doesn't have a permission to retrieve aliases of the group(%s).", user, group))
	}

	aliasEntities := s.repo.ListByGroup(group, time.Now().Add(-time.Hour*24*7))
	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (s AliasServiceImpl) ListByGroupAndAlias(ctx context.Context, user, group, alias string, sort Sort) ([]*domain.Alias, error) {
	if !s.authService.IsInGroup(ctx, user, group) && !s.authService.IsAdmin(ctx, user) {
		return nil, errors.Wrap(ErrNoPermission, fmt.Sprintf("user(%s) doesn't have a permission to retrieve an alias of the group(%s).", user, group))
	}

	if group == "" && alias == "" {
		return s.List(ctx, user)
	}

	aliasEntities := s.repo.ListByGroupAndAlias(group, alias)
	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (s AliasServiceImpl) GoTo(ctx context.Context, user, group, alias string) ([]*domain.Alias, error) {
	if !s.authService.IsInGroup(ctx, user, group) && !s.authService.IsAdmin(ctx, user) {
		return nil, errors.Wrap(ErrNoPermission, fmt.Sprintf("user(%s) doesn't have a permission to retrieve an alias of the group(%s).", user, group))
	}

	aliasEntities := s.repo.ListByGroupAndAlias(group, alias)

	aliases := mapper.AliasesFromEntityToDomain(aliasEntities)
	if err := s.setRecentHits(aliases); err != nil {
		return nil, err
	}
	hit := false
	if len(aliases) >= 1 {
		hit = true
	}
	if s.eventAliasHitRepo == nil {
		log.Warn("eventAliasHitRepo is nil. It might not have been developed yet, so it's just skipped.")
	} else {
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
	}

	return aliases, nil
}

func (s AliasServiceImpl) Update(ctx context.Context, alias *entity.Alias) (*domain.Alias, error) {
	panic("implement me")
}

func (s AliasServiceImpl) Delete(ctx context.Context, id int) (*domain.Alias, error) {
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

func NewAliasService(repo repository.AliasRepository, eventAliasHitRepo repository.EventAliasHitRepository, authService AuthService) AliasService {
	return &AliasServiceImpl{
		repo:              repo,
		eventAliasHitRepo: eventAliasHitRepo,
		authService:       authService,
	}
}

func (s AliasServiceImpl) setRecentHits(aliases []*domain.Alias) error {
	if s.eventAliasHitRepo == nil {
		log.Warn("eventAliasHitRepo is nil. It might not have been developed yet, so it's just skipped.")
		return nil
	}

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

package mapper

import (
	"gorm.io/gorm"
	"rd/domain"
	"rd/entity"
)

func AliasFromDomainToEntity(d *domain.Alias) *entity.Alias {
	return &entity.Alias{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
			DeletedAt: gorm.DeletedAt{
				Time:  d.DeletedAt,
				Valid: true,
			},
		},
		AliasGroup:  d.AliasGroup,
		Name:        d.Name,
		Destination: d.Destination,
	}
}

func AliasFromEntityToDomain(e *entity.Alias) *domain.Alias {
	return &domain.Alias{
		ID:              e.ID,
		AliasGroup:      e.AliasGroup,
		Name:            e.Name,
		Destination:     e.Destination,
		RecentHitCounts: map[string]int{},
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		DeletedAt:       e.DeletedAt.Time,
	}
}
func AliasesFromEntityToDomain(es []*entity.Alias) []*domain.Alias {
	var ds []*domain.Alias
	for _, e := range es {
		ds = append(ds, AliasFromEntityToDomain(e))
	}

	return ds
}

func AliasesFromDomainToEntity(ds []*domain.Alias) []*entity.Alias {
	var es []*entity.Alias
	for _, d := range ds {
		es = append(es, AliasFromDomainToEntity(d))
	}

	return es
}

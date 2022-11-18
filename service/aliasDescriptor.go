package service

import "mytot/domain"

type AliasDescriptorRepository interface {
	List() []*domain.AliasDescriptor
	ListByAlias(alias string) []*domain.AliasDescriptor
}

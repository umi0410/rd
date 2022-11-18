package service

import "rd/domain"

type AliasDescriptorRepository interface {
	List() []*domain.AliasDescriptor
	ListByAlias(alias string) []*domain.AliasDescriptor
}

// +build DEV

package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type ListSourcesDev interface {
	Do() []domain.Source
}

type listSourcesDev struct {
	SourceRepo core.SourceRepository
}

func NewListSourcesDev(Sources core.SourceRepository) ListSourcesDev {
	return &listSourcesDev{SourceRepo: Sources}
}

func (this listSourcesDev) Do() []domain.Source {
	return this.SourceRepo.All()
}

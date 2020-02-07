package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListSourcesDev interface {
	Do() []domain.Source
}

type listSourcesDev struct {
	sourceRepo core.SourceRepository
}

func NewListSourcesDev(Sources core.SourceRepository) ListSourcesDev {
	return &listSourcesDev{sourceRepo: Sources}
}

func (usecase listSourcesDev) Do() []domain.Source {
	return usecase.sourceRepo.All()
}

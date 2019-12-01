package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListStepsDev interface {
	Do() []domain.Step
}

type listStepsDev struct {
	StepRepo core.StepRepository
}

func NewListStepsDev(Steps core.StepRepository) ListStepsDev {
	return &listStepsDev{StepRepo: Steps}
}

func (this listStepsDev) Do() []domain.Step {
	return this.StepRepo.All()
}

package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListStepsDev interface {
	Do() []domain.Step
}

type listStepsDev struct {
	stepRepo core.StepRepository
}

func NewListStepsDev(Steps core.StepRepository) ListStepsDev {
	return &listStepsDev{stepRepo: Steps}
}

func (usecase listStepsDev) Do() []domain.Step {
	return usecase.stepRepo.All()
}

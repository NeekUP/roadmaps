package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListPlansDev interface {
	Do() []domain.Plan
}

type listPlansDev struct {
	planRepo core.PlanRepository
}

func NewListPlansDev(plans core.PlanRepository) ListPlansDev {
	return &listPlansDev{planRepo: plans}
}

func (usecase listPlansDev) Do() []domain.Plan {
	return usecase.planRepo.All()
}

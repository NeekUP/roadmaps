// +build DEV

package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type ListPlansDev interface {
	Do() []domain.Plan
}

type listPlansDev struct {
	PlanRepo core.PlanRepository
}

func NewListPlansDev(plans core.PlanRepository) ListPlansDev {
	return &listPlansDev{PlanRepo: plans}
}

func (this listPlansDev) Do() []domain.Plan {
	return this.PlanRepo.All()
}

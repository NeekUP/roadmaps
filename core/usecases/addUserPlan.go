package usecases

import (
	"roadmaps/core"
)

type AddUserPlan interface {
	Do(ctx core.ReqContext, planId int) (bool, error)
}

type addUserPlan struct {
	PlanRepo      core.PlanRepository
	UsersPlanRepo core.UsersPlanRepository
	Log           core.AppLogger
}

func NewAddUserPlan(planRepo core.PlanRepository, userPlanRepo core.UsersPlanRepository, log core.AppLogger) AddUserPlan {
	return &addUserPlan{
		PlanRepo:      planRepo,
		UsersPlanRepo: userPlanRepo,
		Log:           log,
	}
}

func (this *addUserPlan) Do(ctx core.ReqContext, planId int) (bool, error) {
	plan := this.PlanRepo.Get(planId)
	if plan == nil {
		return false, core.NewError(core.InvalidRequest)
	}

	userId := ctx.UserId()
	success := this.UsersPlanRepo.Add(userId, plan.TopicName, planId)
	if !success {
		return false, core.NewError(core.InternalError)
	}

	return true, nil
}

package usecases

import (
	"github.com/NeekUP/roadmaps/core"
)

type AddUserPlan interface {
	Do(ctx core.ReqContext, planId int) (bool, error)
}

type addUserPlan struct {
	planRepo      core.PlanRepository
	usersPlanRepo core.UsersPlanRepository
	log           core.AppLogger
}

func NewAddUserPlan(planRepo core.PlanRepository, userPlanRepo core.UsersPlanRepository, log core.AppLogger) AddUserPlan {
	return &addUserPlan{
		planRepo:      planRepo,
		usersPlanRepo: userPlanRepo,
		log:           log,
	}
}

func (usecase *addUserPlan) Do(ctx core.ReqContext, planId int) (bool, error) {
	trace := ctx.StartTrace("addUserPlan")
	defer ctx.StopTrace(trace)

	plan := usecase.planRepo.Get(ctx, planId)
	if plan == nil {
		return false, core.NewError(core.InvalidRequest)
	}

	userId := ctx.UserId()
	success, err := usecase.usersPlanRepo.Add(ctx, userId, plan.TopicName, planId)
	if err != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", err.Error(),
		)
		return false, err
	}

	return success, nil
}

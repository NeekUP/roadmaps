package usecases

import "github.com/NeekUP/roadmaps/core"

type RemoveUserPlan interface {
	Do(ctx core.ReqContext, planId int) (bool, error)
}

type removeUserPlan struct {
	usersPlanRepo core.UsersPlanRepository
	log           core.AppLogger
}

func NewRemoveUserPlan(planRepo core.UsersPlanRepository, log core.AppLogger) RemoveUserPlan {
	return &removeUserPlan{
		usersPlanRepo: planRepo,
		log:           log,
	}
}

func (usecase *removeUserPlan) Do(ctx core.ReqContext, planId int) (bool, error) {
	trace := ctx.StartTrace("removeUserPlan")
	defer ctx.StopTrace(trace)

	userId := ctx.UserId()
	if _, err := usecase.usersPlanRepo.Remove(ctx, userId, planId); err != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return false, err
	}
	return true, nil
}

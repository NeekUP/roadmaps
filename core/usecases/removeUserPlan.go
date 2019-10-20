package usecases

import "roadmaps/core"

type RemoveUserPlan interface {
	Do(ctx core.ReqContext, planId int) (bool, error)
}

type removeUserPlan struct {
	UsersPlanRepo core.UsersPlanRepository
	Log           core.AppLogger
}

func NewRemoveUserPlan(planRepo core.UsersPlanRepository, log core.AppLogger) RemoveUserPlan {
	return &removeUserPlan{
		UsersPlanRepo: planRepo,
		Log:           log,
	}
}

func (this *removeUserPlan) Do(ctx core.ReqContext, planId int) (bool, error) {
	userId := ctx.UserId()
	return this.UsersPlanRepo.Remove(userId, planId), nil
}

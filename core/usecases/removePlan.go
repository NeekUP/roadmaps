package usecases

import (
	"github.com/NeekUP/roadmaps/core"
)

type RemovePlan interface {
	Do(ctx core.ReqContext, id int) (bool, error)
}

type removePlan struct {
	PlanRepo core.PlanRepository
	Log      core.AppLogger
}

func NewRemovePlan(planRepo core.PlanRepository, log core.AppLogger) RemovePlan {
	return &removePlan{PlanRepo: planRepo, Log: log}
}

func (this *removePlan) Do(ctx core.ReqContext, id int) (bool, error) {
	appErr := this.validate(id, ctx.UserId())
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	_, err := this.PlanRepo.Delete(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *removePlan) validate(id int, userid string) *core.AppError {
	errors := make(map[string]string)
	if id <= 0 {
		errors["id"] = core.InvalidFormat.String()
	}

	plan := this.PlanRepo.Get(id)
	if plan == nil {
		errors["id"] = core.NotExists.String()
	}

	if plan.OwnerId != userid {
		errors["id"] = core.AccessDenied.String()
	}

	return nil
}

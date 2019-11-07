package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type GetPlan interface {
	Do(ctx core.ReqContext, id int) (*domain.Plan, error)
}

func NewGetPlan(plans core.PlanRepository, users core.UserRepository, logger core.AppLogger) GetPlan {
	return &getPlan{
		PlanRepo: plans,
		UserRepo: users,
		Log:      logger,
	}
}

type getPlan struct {
	PlanRepo core.PlanRepository
	UserRepo core.UserRepository
	Log      core.AppLogger
}

func (this *getPlan) Do(ctx core.ReqContext, id int) (*domain.Plan, error) {
	appErr := this.validate(id)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	plan := this.PlanRepo.Get(id)
	if plan != nil {
		plan.Owner = this.UserRepo.Get(plan.OwnerId)
	}
	return plan, nil
}

func (this *getPlan) validate(id int) *core.AppError {
	errors := make(map[string]string)
	if id < 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

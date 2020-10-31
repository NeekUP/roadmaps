package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetPlanListByUser interface {
	Do(ctx core.ReqContext, count int, page int) ([]domain.Plan, error)
}

func NewGetPlanListByUser(plans core.PlanRepository, users core.UserRepository, logger core.AppLogger) GetPlanListByUser {
	return &getPlanListByUser{
		planRepo: plans,
		userRepo: users,
		log:      logger,
	}
}

type getPlanListByUser struct {
	planRepo core.PlanRepository
	userRepo core.UserRepository
	log      core.AppLogger
}

func (usecase *getPlanListByUser) Do(ctx core.ReqContext, count int, page int) ([]domain.Plan, error) {
	trace := ctx.StartTrace("getPlanListByUser")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(count, page)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"userid", ctx.UserId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}
	userId := ctx.UserId()
	list := usecase.planRepo.GetByUser(ctx, userId, count, page)

	for i := 0; i < len(list); i++ {
		list[i].Owner = usecase.userRepo.Get(ctx, list[i].OwnerId)
	}
	return list, nil
}

func (usecase *getPlanListByUser) validate(count int, page int) *core.AppError {
	errors := make(map[string]string)
	if count <= 0 {
		errors["count"] = core.InvalidCount.String()
	}
	if page < 0 {
		errors["page"] = core.InvalidValue.String()
	}
	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

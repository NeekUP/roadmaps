package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetPlanList interface {
	Do(ctx core.ReqContext, topicName string, count int) ([]domain.Plan, error)
}

func NewGetPlanList(plans core.PlanRepository, users core.UserRepository, logger core.AppLogger) GetPlanList {
	return &getPlanList{
		planRepo: plans,
		userRepo: users,
		log:      logger,
	}
}

type getPlanList struct {
	planRepo core.PlanRepository
	userRepo core.UserRepository
	log      core.AppLogger
}

func (usecase *getPlanList) Do(ctx core.ReqContext, topicName string, count int) ([]domain.Plan, error) {
	trace := ctx.StartTrace("getPlanList")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(topicName, count)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	list := usecase.planRepo.GetPopularByTopic(ctx, topicName, count)

	for i := 0; i < len(list); i++ {
		list[i].Owner = usecase.userRepo.Get(ctx, list[i].OwnerId)
	}
	return list, nil
}

func (usecase *getPlanList) validate(topicName string, count int) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicName(topicName) {
		errors["topicName"] = core.InvalidFormat.String()
	}

	if count <= 0 {
		errors["count"] = core.InvalidCount.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type GetPlanList interface {
	Do(ctx core.ReqContext, topicName string, count int) ([]domain.Plan, error)
}

func NewGetPlanList(plans core.PlanRepository, users core.UserRepository, logger core.AppLogger) GetPlanList {
	return &getPlanList{
		PlanRepo: plans,
		UserRepo: users,
		Log:      logger,
	}
}

type getPlanList struct {
	PlanRepo core.PlanRepository
	UserRepo core.UserRepository
	Log      core.AppLogger
}

func (this *getPlanList) Do(ctx core.ReqContext, topicName string, count int) ([]domain.Plan, error) {
	appErr := this.validate(topicName, count)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	list := this.PlanRepo.GetTopByTopicName(topicName, count)

	for i := 0; i < len(list); i++ {
		list[i].Owner = this.UserRepo.Get(list[i].OwnerId)
	}
	return list, nil
}

func (this *getPlanList) validate(topicName string, count int) *core.AppError {
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

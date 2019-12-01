package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetUsersPlan interface {
	Do(ctx core.ReqContext, topicName string) (*domain.Plan, error)
}

func NewGetUsersPlans(plans core.PlanRepository, usersPlans core.UsersPlanRepository, logger core.AppLogger) GetUsersPlan {
	return &getUsersPlan{
		PlanRepo:   plans,
		UsersPlans: usersPlans,
		Log:        logger,
	}
}

type getUsersPlan struct {
	UsersPlans core.UsersPlanRepository
	PlanRepo   core.PlanRepository
	Log        core.AppLogger
}

func (this *getUsersPlan) Do(ctx core.ReqContext, topicName string) (*domain.Plan, error) {
	appErr := this.validate(topicName)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	userId := ctx.UserId()
	userPlans := this.UsersPlans.GetByTopic(userId, topicName)
	if userPlans == nil {
		return nil, nil
	}

	userPlan := this.PlanRepo.Get(userPlans.PlanId)
	return userPlan, nil
}

func (this *getUsersPlan) validate(topicName string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicName(topicName) {
		errors["name"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

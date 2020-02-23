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
		planRepo:   plans,
		usersPlans: usersPlans,
		log:        logger,
	}
}

type getUsersPlan struct {
	usersPlans core.UsersPlanRepository
	planRepo   core.PlanRepository
	log        core.AppLogger
}

func (usecase *getUsersPlan) Do(ctx core.ReqContext, topicName string) (*domain.Plan, error) {
	trace := ctx.StartTrace("getUsersPlan")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(topicName)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	userId := ctx.UserId()
	userPlans := usecase.usersPlans.GetByTopic(ctx, userId, topicName)
	if userPlans == nil {
		return nil, nil
	}

	userPlan := usecase.planRepo.Get(ctx, userPlans.PlanId)
	return userPlan, nil
}

func (usecase *getUsersPlan) validate(topicName string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicName(topicName) {
		errors["name"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

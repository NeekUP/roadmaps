package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetTopic interface {
	Do(ctx core.ReqContext, name string, planCount int) (*domain.Topic, error)
	//DoById(ctx core.ReqContext, id int, planCount int) (*domain.Topic, error)
}

type getTopic struct {
	topicRepo core.TopicRepository
	planRepo  core.PlanRepository
	usersPlan core.UsersPlanRepository
	log       core.AppLogger
}

func NewGetTopic(topicRepo core.TopicRepository, planRepo core.PlanRepository, userPlans core.UsersPlanRepository, log core.AppLogger) GetTopic {
	return &getTopic{topicRepo: topicRepo, planRepo: planRepo, usersPlan: userPlans, log: log}
}

func (usecase *getTopic) Do(ctx core.ReqContext, name string, planCount int) (*domain.Topic, error) {
	trace := ctx.StartTrace("getTopic")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(name, planCount)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	topic := usecase.topicRepo.Get(ctx, name)
	if topic == nil {
		return nil, core.NewError(core.NotExists)
	}
	return topic, nil
}

func (usecase *getTopic) validate(name string, planCount int) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicName(name) {
		errors["name"] = core.InvalidFormat.String()
	}

	if planCount <= 0 {
		errors["count"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

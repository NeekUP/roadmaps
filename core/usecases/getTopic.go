package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type GetTopic interface {
	Do(ctx core.ReqContext, name string, planCount int) (*domain.Topic, error)
	DoById(ctx core.ReqContext, id int, planCount int) (*domain.Topic, error)
}

type getTopic struct {
	TopicRepo core.TopicRepository
	PlanRepo  core.PlanRepository
	UsersPlan core.UsersPlanRepository
	Log       core.AppLogger
}

func NewGetTopic(topicRepo core.TopicRepository, planRepo core.PlanRepository, userPlans core.UsersPlanRepository, log core.AppLogger) GetTopic {
	return &getTopic{TopicRepo: topicRepo, PlanRepo: planRepo, UsersPlan: userPlans, Log: log}
}

func (this *getTopic) Do(ctx core.ReqContext, name string, planCount int) (*domain.Topic, error) {
	appErr := this.validate(name, planCount)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	topic := this.TopicRepo.Get(name)
	if topic == nil {
		return nil, core.NewError(core.NotExists)
	}

	return topic, nil
}

func (this *getTopic) DoById(ctx core.ReqContext, id int, planCount int) (*domain.Topic, error) {

	topic := this.TopicRepo.GetById(id)
	if topic == nil {
		return nil, core.NewError(core.NotExists)
	}

	if len(topic.Plans) >= planCount {
		return topic, nil
	}

	return topic, nil
}

// func (this *getTopic) AttachePlans(ctx core.ReqContext, topic *domain.Topic, planCount int, includeSteps bool) {
// 	if ctx.UserId() != "" {
// 		upId := this.UsersPlan.GetByTopic(ctx.UserId(), topic.Name)
// 		if upId != nil {
// 			userSelectedPlan := this.PlanRepo.Get(upId.PlanId)
// 			if userSelectedPlan != nil {
// 				topic.Plans = append(topic.Plans, *userSelectedPlan)
// 			}
// 		}
// 	}

// 	if len(topic.Plans) == 0 {
// 		topic.Plans = this.PlanRepo.GetTopByTopicName(topic.Name, planCount, includeSteps)
// 	} else if planCount-len(topic.Plans) == 0 {
// 		return
// 	} else {
// 		u := topic.Plans[0]
// 		topic.Plans = this.PlanRepo.GetTopByTopicName(topic.Name, planCount-len(topic.Plans), includeSteps, topic.Plans[0].Id)
// 		topic.Plans = append(topic.Plans, u)
// 	}
// }

func (this *getTopic) validate(name string, planCount int) *core.AppError {
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

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
	Log       core.AppLogger
}

func NewGetTopic(topicRepo core.TopicRepository, planRepo core.PlanRepository, log core.AppLogger) GetTopic {
	return &getTopic{TopicRepo: topicRepo, PlanRepo: planRepo, Log: log}
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

	// TODO: Find user selected plan. if exists

	//
	if len(topic.Plans) >= planCount {
		return topic, nil
	}

	topic.Plans = this.PlanRepo.GetTopByTopicName(topic.Name, planCount)
	return topic, nil
}

func (this *getTopic) DoById(ctx core.ReqContext, id int, planCount int) (*domain.Topic, error) {

	topic := this.TopicRepo.GetById(id)
	if topic == nil {
		return nil, core.NewError(core.NotExists)
	}

	// TODO: Find user selected plan. if exists

	//
	if len(topic.Plans) >= planCount {
		return topic, nil
	}

	topic.Plans = this.PlanRepo.GetTopByTopicName(topic.Name, planCount)
	return topic, nil
}

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

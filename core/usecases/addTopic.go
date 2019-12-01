package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddTopic interface {
	Do(ctx core.ReqContext, title, desc string) (*domain.Topic, error)
}

type addTopic struct {
	TopicRepo core.TopicRepository
	Log       core.AppLogger
}

func NewAddTopic(topicRepo core.TopicRepository, log core.AppLogger) AddTopic {
	return &addTopic{TopicRepo: topicRepo, Log: log}
}

func (this *addTopic) Do(ctx core.ReqContext, title, desc string) (*domain.Topic, error) {
	appErr := this.validate(title, desc)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	userId := ctx.UserId()
	topic := domain.NewTopic(title, desc, userId)

	saved, err := this.TopicRepo.Save(topic)
	if err != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return nil, err
	}

	if saved {
		return topic, nil
	}

	return nil, nil
}

func (this *addTopic) validate(title, desc string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicTitle(title) {
		errors["title"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

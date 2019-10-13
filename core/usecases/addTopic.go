package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
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

	// check retrictions in db
	saved := this.TopicRepo.Save(topic)
	if saved {
		return topic, nil
	}

	return nil, core.NewError(core.AlreadyExists)
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
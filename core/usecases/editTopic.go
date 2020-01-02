package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type EditTopic interface {
	Do(ctx core.ReqContext, id int, title, desc string, istag bool) (bool, error)
}

type editTopic struct {
	TopicRepo core.TopicRepository
	Log       core.AppLogger
}

func NewEditTopic(topicRepo core.TopicRepository, log core.AppLogger) EditTopic {
	return &editTopic{TopicRepo: topicRepo, Log: log}
}

func (et *editTopic) Do(ctx core.ReqContext, id int, title, desc string, istag bool) (bool, error) {
	appErr := et.validate(id, title)
	if appErr != nil {
		et.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	userId := ctx.UserId()
	topic := domain.NewTopic(title, desc, userId)
	topic.Id = id
	topic.IsTag = istag
	saved, err := et.TopicRepo.Update(topic)
	if err != nil {
		et.Log.Errorw("Topic not updated",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return false, err
	}

	return saved, nil
}

func (et *editTopic) validate(id int, title string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicTitle(title) {
		errors["title"] = core.InvalidFormat.String()
	}

	if id <= 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

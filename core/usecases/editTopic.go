package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type EditTopic interface {
	Do(ctx core.ReqContext, id int, title, desc string, istag bool) (bool, error)
}

type editTopic struct {
	repo      core.TopicRepository
	log       core.AppLogger
	changeLog core.ChangeLog
}

func NewEditTopic(topicRepo core.TopicRepository, changelog core.ChangeLog, log core.AppLogger) EditTopic {
	return &editTopic{repo: topicRepo, changeLog: changelog, log: log}
}

func (usecase *editTopic) Do(ctx core.ReqContext, id int, title, desc string, istag bool) (bool, error) {
	userId := ctx.UserId()
	old := usecase.repo.GetById(id)
	appErr := usecase.validate(id, title, old, userId)

	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	topic := domain.NewTopic(title, desc, userId)
	topic.Id = id
	topic.IsTag = istag
	saved, err := usecase.repo.Update(topic)
	if err != nil {
		usecase.log.Errorw("Topic not updated",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return false, err
	}

	usecase.changeLog.Edited(domain.TopicEntity, int64(topic.Id), ctx.UserId(), old, topic)
	return saved, nil
}

func (usecase *editTopic) validate(id int, title string, topic *domain.Topic, userId string) *core.AppError {
	errors := make(map[string]string)

	if topic == nil {
		errors["id"] = core.NotExists.String()
	}

	if topic.Creator != userId {
		errors["id"] = core.AccessDenied.String()
	}

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

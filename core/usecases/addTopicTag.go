package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddTopicTag interface {
	Do(ctx core.ReqContext, tagname, topicname string) (bool, error)
}

type addTopicTag struct {
	topicRepo core.TopicRepository
	log       core.AppLogger
	changeLog core.ChangeLog
}

func NewAddTopicTag(topicRepo core.TopicRepository, changeLog core.ChangeLog, log core.AppLogger) AddTopicTag {
	return &addTopicTag{topicRepo: topicRepo, changeLog: changeLog, log: log}
}

func (usecase *addTopicTag) Do(ctx core.ReqContext, tagname, topicname string) (bool, error) {
	trace := ctx.StartTrace("addTopicTag")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(tagname, topicname)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	userId := ctx.UserId()
	topic := usecase.topicRepo.Get(ctx, topicname)
	if topic == nil {
		return false, core.NewError(core.NotExists)
	}

	hasChanges := usecase.topicRepo.AddTag(ctx, tagname, topicname)
	if hasChanges {
		changedTopic := *topic
		copy(changedTopic.Tags, topic.Tags)
		changedTopic.Tags = append(changedTopic.Tags, domain.TopicTag{Name: topicname})
		usecase.changeLog.Edited(domain.TopicEntity, int64(topic.Id), userId, topic, &changedTopic)
	}
	return hasChanges, nil
}

func (usecase *addTopicTag) validate(tagname, topicname string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicName(tagname) {
		errors["tagname"] = core.InvalidFormat.String()
	}

	if !core.IsValidTopicName(topicname) {
		errors["topicname"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

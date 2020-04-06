package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type RemoveTopicTag interface {
	Do(ctx core.ReqContext, tagname, topicname string) (bool, error)
}

type removeTopicTag struct {
	topicRepo core.TopicRepository
	log       core.AppLogger
	changeLog core.ChangeLog
}

func NewRemoveTopicTag(topicRepo core.TopicRepository, changeLog core.ChangeLog, log core.AppLogger) RemoveTopicTag {
	return &removeTopicTag{topicRepo: topicRepo, changeLog: changeLog, log: log}
}

func (usecase *removeTopicTag) Do(ctx core.ReqContext, tagname, topicname string) (bool, error) {
	trace := ctx.StartTrace("removeTopicTag")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(tagname, topicname)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return false, appErr
	}

	userId := ctx.UserId()
	topic := usecase.topicRepo.Get(ctx, topicname)
	if topic == nil {
		return false, core.NewError(core.NotExists)
	}

	result := usecase.topicRepo.DeleteTag(ctx, tagname, topicname)
	if result {
		changedTopic := *topic
		if len(topic.Tags) > 1 {
			changedTopic.Tags = make([]domain.TopicTag, len(topic.Tags)-1)
			p := 0
			for i, _ := range changedTopic.Tags {
				if topic.Tags[i+p].Name == topicname {
					p++
				}
				changedTopic.Tags[i] = topic.Tags[i+p]
			}
		} else {
			changedTopic.Tags = make([]domain.TopicTag, 0)
		}

		usecase.changeLog.Edited(domain.TopicEntity, int64(topic.Id), userId, topic, &changedTopic)
	}
	return result, nil
}

func (usecase *removeTopicTag) validate(tagname, topicname string) *core.AppError {
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

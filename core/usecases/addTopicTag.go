package usecases

import (
	"github.com/NeekUP/roadmaps/core"
)

type AddTopicTag interface {
	Do(ctx core.ReqContext, tagname, topicname string) (bool, error)
}

type addTopicTag struct {
	TopicRepo core.TopicRepository
	Log       core.AppLogger
}

func NewAddTopicTag(topicRepo core.TopicRepository, log core.AppLogger) AddTopicTag {
	return &addTopicTag{TopicRepo: topicRepo, Log: log}
}

func (a *addTopicTag) Do(ctx core.ReqContext, tagname, topicname string) (bool, error) {
	appErr := a.validate(tagname, topicname)
	if appErr != nil {
		a.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	result := a.TopicRepo.AddTag(tagname, topicname)
	return result, nil
}

func (a *addTopicTag) validate(tagname, topicname string) *core.AppError {
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

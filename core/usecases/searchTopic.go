package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type SearchTopic interface {
	Do(ctx core.ReqContext, str string, count int) []domain.Topic
}

type searchTopic struct {
	TopicRepo core.TopicRepository
	Log       core.AppLogger
}

func NewSearchTopic(topicRepo core.TopicRepository, log core.AppLogger) SearchTopic {
	return &searchTopic{TopicRepo: topicRepo, Log: log}
}

func (r *searchTopic) Do(ctx core.ReqContext, str string, count int) []domain.Topic {
	appErr := r.validate(str)
	if appErr != nil {
		r.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return []domain.Topic{}
	}

	return r.TopicRepo.TitleLike(str, count)
}

func (r *searchTopic) validate(str string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicTitle(str) {
		errors["search"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

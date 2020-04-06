package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type SearchTopic interface {
	Do(ctx core.ReqContext, str string, count int) []domain.Topic
}

type searchTopic struct {
	topicRepo core.TopicRepository
	log       core.AppLogger
}

func NewSearchTopic(topicRepo core.TopicRepository, log core.AppLogger) SearchTopic {
	return &searchTopic{topicRepo: topicRepo, log: log}
}

func (usecase *searchTopic) Do(ctx core.ReqContext, str string, count int) []domain.Topic {
	trace := ctx.StartTrace("searchTopic")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(str)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return []domain.Topic{}
	}

	return usecase.topicRepo.TitleLike(ctx, str, count)
}

func (usecase *searchTopic) validate(str string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicTitle(str) {
		errors["search"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

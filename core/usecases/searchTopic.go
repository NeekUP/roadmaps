package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

const MAX_TOPIC_SEARCH_RESULTS_SOUNT = 30

type SearchTopic interface {
	Do(ctx core.ReqContext, str string, tags []string, count int) []domain.Topic
}

type searchTopic struct {
	topicRepo core.TopicRepository
	log       core.AppLogger
}

func NewSearchTopic(topicRepo core.TopicRepository, log core.AppLogger) SearchTopic {
	return &searchTopic{topicRepo: topicRepo, log: log}
}

func (usecase *searchTopic) Do(ctx core.ReqContext, str string, tags []string, count int) []domain.Topic {
	trace := ctx.StartTrace("searchTopic")
	defer ctx.StopTrace(trace)

	count = usecase.adjustResultsCount(count)

	appErr := usecase.validate(str)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return []domain.Topic{}
	}

	return usecase.topicRepo.Search(ctx, str, tags, count)
}

func (usecase *searchTopic) adjustResultsCount(count int) int {
	if count == 0 {
		count = 10
	} else if count > MAX_TOPIC_SEARCH_RESULTS_SOUNT {
		count = MAX_TOPIC_SEARCH_RESULTS_SOUNT
	}
	return count
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

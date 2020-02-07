package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddTopic interface {
	Do(ctx core.ReqContext, title, desc string, istag bool, tags []string) (*domain.Topic, error)
}

type addTopic struct {
	topicRepo core.TopicRepository
	log       core.AppLogger
	changeLog core.ChangeLog
}

func NewAddTopic(topicRepo core.TopicRepository, changelog core.ChangeLog, log core.AppLogger) AddTopic {
	return &addTopic{topicRepo: topicRepo, changeLog: changelog, log: log}
}

func (usecase *addTopic) Do(ctx core.ReqContext, title, desc string, istag bool, tags []string) (*domain.Topic, error) {
	appErr := usecase.validate(title, tags)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	userId := ctx.UserId()
	topic := domain.NewTopic(title, desc, userId)
	topic.IsTag = istag
	topic.Tags = usecase.topicRepo.GetTags(tags)
	saved, err := usecase.topicRepo.Save(topic)
	if err != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return nil, err
	}

	if saved {
		if len(topic.Tags) > 0 {
			for _, tag := range topic.Tags {
				usecase.topicRepo.AddTag(tag.Name, topic.Name)
			}
		}
		usecase.changeLog.Added(domain.TopicEntity, int64(topic.Id), userId)
		return topic, nil
	}

	return nil, nil
}

func (usecase *addTopic) validate(title string, tags []string) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicTitle(title) {
		errors["title"] = core.InvalidFormat.String()
	}

	for _, tag := range tags {
		if !core.IsValidTopicName(tag) {
			errors["tags"] = core.InvalidFormat.String()
			break
		}
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

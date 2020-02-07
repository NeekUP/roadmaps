package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListTopicsDev interface {
	Do() []domain.Topic
}

type listTopicsDev struct {
	topicRepo core.TopicRepository
}

func NewListTopicsDev(topics core.TopicRepository) ListTopicsDev {
	return &listTopicsDev{topicRepo: topics}
}

func (usecase listTopicsDev) Do() []domain.Topic {
	return usecase.topicRepo.All()
}

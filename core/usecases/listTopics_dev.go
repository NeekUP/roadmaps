package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListTopicsDev interface {
	Do() []domain.Topic
}

type listTopicsDev struct {
	TopicRepo core.TopicRepository
}

func NewListTopicsDev(topics core.TopicRepository) ListTopicsDev {
	return &listTopicsDev{TopicRepo: topics}
}

func (this listTopicsDev) Do() []domain.Topic {
	return this.TopicRepo.All()
}

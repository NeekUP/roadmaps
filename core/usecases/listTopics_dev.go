// +build DEV

package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
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

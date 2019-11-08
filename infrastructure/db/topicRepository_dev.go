// +build DEV

package db

import (
	"database/sql"
	"roadmaps/core"
	"roadmaps/domain"
	"sync"
)

var (
	Topics        = make([]domain.Topic, 0)
	TopicsMux     sync.Mutex
	TopicsCounter int
)

type topicRepoInMemory struct {
	Conn *sql.DB
}

func NewTopicRepository(conn *sql.DB) core.TopicRepository {
	return &topicRepoInMemory{
		Conn: conn}
}

func (this *topicRepoInMemory) Get(name string) *domain.Topic {
	TopicsMux.Lock()
	defer TopicsMux.Unlock()

	for i := 0; i < len(Topics); i++ {
		if Topics[i].Name == name {
			copy := Topics[i]
			return &copy
		}
	}
	return nil
}

func (this *topicRepoInMemory) GetById(id int) *domain.Topic {
	TopicsMux.Lock()
	defer TopicsMux.Unlock()

	for i := 0; i < len(Topics); i++ {
		if Topics[i].Id == id {
			copy := Topics[i]
			return &copy
		}
	}
	return nil
}

func (this *topicRepoInMemory) Save(topic *domain.Topic) bool {
	TopicsMux.Lock()
	defer TopicsMux.Unlock()

	for i := 0; i < len(Topics); i++ {
		if Topics[i].Name == topic.Name {
			return false
		}
	}

	TopicsCounter++
	topic.Id = TopicsCounter
	Topics = append(Topics, *topic)
	return true
}

func (this *topicRepoInMemory) Update(topic *domain.Topic) bool {
	TopicsMux.Lock()
	defer TopicsMux.Unlock()

	for i := 0; i < len(Topics); i++ {
		if Topics[i].Id == topic.Id {
			Topics[i] = *topic
			return true
		}
	}
	return false
}

//
//
package db

//
//import (
//	"sync"
//
//	"github.com/NeekUP/roadmaps/core"
//	"github.com/NeekUP/roadmaps/domain"
//	"github.com/jackc/pgx/v4"
//)
//
//var (
//	Topics        = make([]domain.TopicName, 0)
//	TopicsMux     sync.Mutex
//	TopicsCounter int
//)
//
//type topicRepoInMemory struct {
//	Conn *pgx.Conn
//}
//
//func NewTopicRepository(conn *pgx.Conn) core.TopicRepository {
//	return &topicRepoInMemory{
//		Conn: conn}
//}
//
//func (this *topicRepoInMemory) Get(name string) *domain.TopicName {
//	TopicsMux.Lock()
//	defer TopicsMux.Unlock()
//
//	for i := 0; i < len(Topics); i++ {
//		if Topics[i].Name == name {
//			copy := Topics[i]
//			return &copy
//		}
//	}
//	return nil
//}
//
//func (this *topicRepoInMemory) GetById(id int) *domain.TopicName {
//	TopicsMux.Lock()
//	defer TopicsMux.Unlock()
//
//	for i := 0; i < len(Topics); i++ {
//		if Topics[i].Id == id {
//			copy := Topics[i]
//			return &copy
//		}
//	}
//	return nil
//}
//
//func (this *topicRepoInMemory) Save(topic *domain.TopicName) bool {
//	TopicsMux.Lock()
//	defer TopicsMux.Unlock()
//
//	for i := 0; i < len(Topics); i++ {
//		if Topics[i].Name == topic.Name {
//			return false
//		}
//	}
//
//	TopicsCounter++
//	topic.Id = TopicsCounter
//	Topics = append(Topics, *topic)
//	return true
//}
//
//func (this *topicRepoInMemory) Update(topic *domain.TopicName) bool {
//	TopicsMux.Lock()
//	defer TopicsMux.Unlock()
//
//	for i := 0; i < len(Topics); i++ {
//		if Topics[i].Id == topic.Id {
//			Topics[i] = *topic
//			return true
//		}
//	}
//	return false
//}
//
//func (this *topicRepoInMemory) All() []domain.TopicName {
//	return Topics
//}

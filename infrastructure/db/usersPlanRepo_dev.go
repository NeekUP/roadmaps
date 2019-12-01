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
//	UsersPlans    = make([]domain.UsersPlan, 0)
//	UsersPlansMux sync.Mutex
//)
//
//type usersPlanRepoInMemory struct {
//	Conn *pgx.Conn
//}
//
//func NewUsersPlanRepository(conn *pgx.Conn) core.UsersPlanRepository {
//	return &usersPlanRepoInMemory{
//		Conn: conn}
//}
//
//func (this *usersPlanRepoInMemory) Add(userId string, topicName string, planId int) bool {
//	UsersPlansMux.Lock()
//	defer UsersPlansMux.Unlock()
//
//	for i := 0; i < len(UsersPlans); i++ {
//		if UsersPlans[i].UserId == userId &&
//			UsersPlans[i].TopicName == topicName {
//			if i == len(UsersPlans)-1 {
//				UsersPlans = UsersPlans[:len(UsersPlans)-1]
//			} else {
//				UsersPlans = append(UsersPlans[:i], UsersPlans[i+1:]...)
//			}
//		}
//	}
//
//	UsersPlans = append(UsersPlans, domain.UsersPlan{UserId: userId, TopicName: topicName, PlanId: planId})
//	return true
//}
//
//func (this *usersPlanRepoInMemory) Remove(userId string, planId int) bool {
//	UsersPlansMux.Lock()
//	defer UsersPlansMux.Unlock()
//
//	l := len(UsersPlans)
//	for i := 0; i < l; i++ {
//		if UsersPlans[i].UserId == userId &&
//			UsersPlans[i].PlanId == planId {
//			if i == l-1 {
//				UsersPlans = UsersPlans[:l-1]
//				return true
//			}
//			UsersPlans = append(UsersPlans[:i], UsersPlans[i+1:]...)
//			return true
//		}
//	}
//	return true
//}
//
//func (this *usersPlanRepoInMemory) GetByTopic(userId, topicName string) *domain.UsersPlan {
//	UsersPlansMux.Lock()
//	defer UsersPlansMux.Unlock()
//
//	l := len(UsersPlans)
//	for i := 0; i < l; i++ {
//		if UsersPlans[i].UserId == userId &&
//			UsersPlans[i].TopicName == topicName {
//			copy := UsersPlans[i]
//			return &copy
//		}
//	}
//	return nil
//}
//
//func (this *usersPlanRepoInMemory) GetByUser(userId string) []domain.UsersPlan {
//	UsersPlansMux.Lock()
//	defer UsersPlansMux.Unlock()
//
//	usersPlans := []domain.UsersPlan{}
//
//	l := len(UsersPlans)
//	for i := 0; i < l; i++ {
//		if UsersPlans[i].UserId == userId {
//			usersPlans = append(usersPlans, UsersPlans[i])
//		}
//	}
//	return usersPlans
//}

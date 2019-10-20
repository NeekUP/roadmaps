// +build DEV

package db

import (
	"database/sql"
	"roadmaps/core"
	"sync"
)

var (
	UsersPlans    = make([]usersPlanRecord, 0)
	UsersPlansMux sync.Mutex
)

type usersPlanRecord struct {
	UserId    string
	PlanId    int
	TopicName string
}

type usersPlanRepoInMemory struct {
	Conn *sql.DB
}

func NewUsersPlanRepository(conn *sql.DB) core.UsersPlanRepository {
	return &usersPlanRepoInMemory{
		Conn: conn}
}

func (this *usersPlanRepoInMemory) Add(userId string, topicName string, planId int) bool {
	UsersPlansMux.Lock()
	defer UsersPlansMux.Unlock()

	for i := 0; i < len(UsersPlans); i++ {
		if UsersPlans[i].UserId == userId &&
			UsersPlans[i].TopicName == topicName {
			UsersPlans[i].PlanId = planId
			return true
		}
	}

	UsersPlans = append(UsersPlans, usersPlanRecord{userId, planId, topicName})
	return true
}

func (this *usersPlanRepoInMemory) Remove(userId string, planId int) bool {
	UsersPlansMux.Lock()
	defer UsersPlansMux.Unlock()

	l := len(UsersPlans)
	for i := 0; i < l; i++ {
		if UsersPlans[i].UserId == userId &&
			UsersPlans[i].PlanId == planId {
			if i == l-1 {
				UsersPlans = UsersPlans[:l-1]
				return true
			}
			UsersPlans = append(UsersPlans[:i], UsersPlans[i+1:]...)
			return true
		}
	}
	return true
}

func (this *usersPlanRepoInMemory) GetByTopic(userId, topicName string) (planId int, exists bool) {
	UsersPlansMux.Lock()
	defer UsersPlansMux.Unlock()

	l := len(UsersPlans)
	for i := 0; i < l; i++ {
		if UsersPlans[i].UserId == userId &&
			UsersPlans[i].TopicName == topicName {
			return UsersPlans[i].PlanId, true
		}
	}
	return 0, false
}

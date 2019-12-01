package db

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type usersPlanRepo struct {
	Db *DbConnection
}

func NewUsersPlanRepository(db *DbConnection) core.UsersPlanRepository {
	return &usersPlanRepo{
		Db: db}
}

func (this *usersPlanRepo) Add(userId string, topicName string, planId int) (bool, *core.AppError) {
	panic("implement me")
}

func (this *usersPlanRepo) Remove(userId string, planId int) (bool, *core.AppError) {
	panic("implement me")
}

func (this *usersPlanRepo) GetByTopic(userId, topicName string) *domain.UsersPlan {
	panic("implement me")
}

func (this *usersPlanRepo) GetByUser(userId string) []domain.UsersPlan {
	panic("implement me")
}

package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListUsersDev interface {
	Do() []domain.User
}

type listUsersDev struct {
	UserRepo core.UserRepository
}

func NewListUsersDev(Users core.UserRepository) ListUsersDev {
	return &listUsersDev{UserRepo: Users}
}

func (this *listUsersDev) Do() []domain.User {
	return this.UserRepo.All()
}

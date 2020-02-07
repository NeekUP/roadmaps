package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type ListUsersDev interface {
	Do() []domain.User
}

type listUsersDev struct {
	userRepo core.UserRepository
}

func NewListUsersDev(Users core.UserRepository) ListUsersDev {
	return &listUsersDev{userRepo: Users}
}

func (usecase *listUsersDev) Do() []domain.User {
	return usecase.userRepo.All()
}

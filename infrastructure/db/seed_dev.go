package db

import (
	"github.com/google/uuid"
	"roadmaps/domain"
)

func SeedUsers() []domain.User {

	users := make([]domain.User, 10)
	users = append(users, domain.User{
		Id:             uuid.New().String(),
		Email:          "nikita@popovsky.pro",
		EmailConfirmed: true,
		Rights:         domain.U | domain.M | domain.A | domain.O})

	return users
}

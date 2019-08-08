// +build DEV

package infrastructure

import (
	"roadmaps/core"
	"roadmaps/core/usecases"
)

type DbSeed interface {
	Seed()
}

func NewDbSeed(regUser usecases.RegisterUser, userRepo core.UserRepository) DbSeed {
	return &dbSeedDev{
		RegUser:  regUser,
		UserRepo: userRepo,
	}
}

type dbSeedDev struct {
	RegUser  usecases.RegisterUser
	UserRepo core.UserRepository
}

func (this *dbSeedDev) Seed() {
	this.seedUsers()
}

func (this *dbSeedDev) seedUsers() {
	context := NewContext(nil)
	if !this.UserRepo.ExistsEmail("nikita@popovsky.pro") {
		this.RegUser.Do(context, "Neek", "nikita@popovsky.pro", "123456")
	}
}

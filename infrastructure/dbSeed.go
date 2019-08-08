// +build PROD

package infrastructure

import (
	"roadmaps/core"
	"roadmaps/core/usecases"
)

func NewDbSeed(regUser usecases.RegisterUser, userRepo core.UserRepository) DbSeed {
	return &dbSeedProd{
		RegUser:  regUser,
		UserRepo: userRepo,
	}
}

type dbSeedProd struct {
	RegUser  usecases.RegisterUser
	UserRepo core.UserRepository
}

func (this *dbSeedProd) Seed() {
	this.seedUsers()
}

func (this *dbSeedProd) seedUsers() {
	context := &fakeContext{}
	if !this.UserRepo.ExistsEmail("nikita@popovsky.pro") {
		this.RegUser.Do(context, "Neek", "nikita@popovsky.pro", "123456")
	}

}

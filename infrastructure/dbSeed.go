// +build PROD

package infrastructure

import (
	"fmt"
	"os"
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
	if this.UserRepo.Count() == 0 {
		context := NewContext(nil)
		for i := 0; i < 10; i++ {
			name := os.Getenv(fmt.Sprintf("adminname%d", i))
			email := os.Getenv(fmt.Sprintf("adminemail%d", i))
			pass := os.Getenv(fmt.Sprintf("adminpass%d", i))

			if !this.UserRepo.ExistsEmail(email) && name != "" && email != "" && pass !=  {
				this.RegUser.Do(context, name, email, pass)
			}
		}
	}
}

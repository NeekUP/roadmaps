package infrastructure

import (
	"context"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"os"
	"strings"
)

type DbSeed interface {
	Seed()
}

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

func (seed *dbSeedProd) Seed() {
	seed.seedUsers()
}

func (seed *dbSeedProd) seedUsers() {
	if count, ok := seed.UserRepo.Count(); ok && count == 0 {
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			fmt.Println(pair[0])
		}
		for i := 0; i < 10; i++ {
			name := os.Getenv(fmt.Sprintf("adminname%d", i))
			email := os.Getenv(fmt.Sprintf("adminemail%d", i))
			pass := os.Getenv(fmt.Sprintf("adminpass%d", i))

			exists, ok := seed.UserRepo.ExistsEmail(email)
			if !exists && ok && name != "" && email != "" && pass != "" {
				seed.RegUser.Do(NewContext(context.Background()), name, email, pass)
			}
		}
	}
}

// +build DEV

package infrastructure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/domain"
	"roadmaps/infrastructure/db"
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
	this.seedEntities()
}

func (this *dbSeedDev) seedUsers() {
	context := NewContext(nil)
	if !this.UserRepo.ExistsEmail("neek@neek.com") {
		this.RegUser.Do(context, "Neek", "neek@neek.com", "123456")
	}
	if !this.UserRepo.ExistsEmail("alen@alen.com") {
		this.RegUser.Do(context, "Alen", "alen@alen.com", "123456")
	}
}

func (this *dbSeedDev) seedEntities() {

	users := make(map[string]*domain.User)
	users["Neek"] = this.UserRepo.FindByEmail("neek@neek.com")
	users["Alen"] = this.UserRepo.FindByEmail("alen@alen.com")

	jsonFile, err := os.Open("dev_db.json")
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var devdb DevDB
	err = json.Unmarshal(byteValue, &devdb)
	if err != nil {
		panic("fail to unmarshal dev_db.json content")
	}

	db.SourcesCounter = len(devdb.Resources)
	db.Sources = make([]domain.Source, db.SourcesCounter)
	for i, v := range devdb.Resources {
		db.Sources[i] = v
	}

	db.TopicsCounter = len(devdb.Topics)
	db.Topics = make([]domain.Topic, db.TopicsCounter)
	for i, v := range devdb.Topics {
		v.Creator = users[v.Creator].Id
		db.Topics[i] = v
	}

	db.PlansCounter = len(devdb.Plans)
	db.Plans = make([]domain.Plan, db.PlansCounter)
	for i, v := range devdb.Plans {
		v.OwnerId = users[v.OwnerId].Id
		db.Plans[i] = v
	}

	db.StepsCounter = len(devdb.Steps)
	db.Steps = make([]domain.Step, db.StepsCounter)
	for i, v := range devdb.Steps {
		db.Steps[i] = v
	}

	db.UsersPlans = make([]domain.UsersPlan, len(devdb.UsersPlans))
	for i, v := range devdb.UsersPlans {
		v.UserId = users[v.UserId].Id
		db.UsersPlans[i] = v
	}
}

type DevDB struct {
	Resources  []domain.Source    `json:"resources"`
	Topics     []domain.Topic     `json:"topics"`
	Plans      []domain.Plan      `json:"plans"`
	Steps      []domain.Step      `json:"steps"`
	UsersPlans []domain.UsersPlan `json:"usersPlans"`
}

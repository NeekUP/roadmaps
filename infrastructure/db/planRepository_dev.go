// +build DEV

package db

import (
	"database/sql"
	"fmt"
	"roadmaps/core"
	"roadmaps/domain"
	"sort"
	"sync"
)

var (
	Plans        = make([]domain.Plan, 0)
	PlansMux     sync.Mutex
	PlansCounter int
)

type planRepoInMemory struct {
	Conn *sql.DB
}

func NewPlansRepository(conn *sql.DB) core.PlanRepository {
	return &planRepoInMemory{
		Conn: conn}
}

func (this *planRepoInMemory) SaveWithSteps(plan *domain.Plan) error {
	TopicsMux.Lock()
	defer TopicsMux.Unlock()

	var topic *domain.Topic
	for i := 0; i < len(Topics); i++ {
		if Topics[i].Name == plan.TopicName {
			topic = &Topics[i]
		}
	}

	if topic == nil {
		return fmt.Errorf("Topic not found be name: %s", plan.TopicName)
	}

	if !this.Save(plan) {
		return fmt.Errorf("Plan with same title already exists: %s", plan.Title)
	}

	StepsMux.Lock()
	defer StepsMux.Unlock()

	for i := 0; i < len(plan.Steps); i++ {
		StepsCounter++
		plan.Steps[i].PlanId = plan.Id
		plan.Steps[i].Id = StepsCounter
		Steps = append(Steps, plan.Steps[i])
	}

	return nil
}

func (this *planRepoInMemory) Get(id int) *domain.Plan {
	PlansMux.Lock()
	defer PlansMux.Unlock()

	for _, v := range Plans {
		if v.Id == id {
			result := v

			for _, s := range Steps {
				if s.PlanId == id {
					result.Steps = append(result.Steps, s)
				}
			}

			return &result
		}
	}

	return nil
}

func (this *planRepoInMemory) GetList(id []int) []domain.Plan {
	list := make([]domain.Plan, 0)
	for i := 0; i < len(id); i++ {
		p := this.Get(id[i])
		if p != nil {
			list = append(list, *p)
		}
	}

	return list
}

func (this *planRepoInMemory) GetTopByTopicName(topic string, count int) []domain.Plan {

	PlansMux.Lock()
	defer PlansMux.Unlock()

	topicPlans := make([]domain.Plan, 0)

	for i := 0; i < PlansCounter; i++ {
		if Plans[i].TopicName == topic {
			topicPlans = append(topicPlans, Plans[i])
		}
	}

	sort.Slice(topicPlans, func(i, j int) bool {
		return topicPlans[i].Points > topicPlans[j].Points
	})

	if count > len(topicPlans) {
		count = len(topicPlans)
	}

	return topicPlans[:count]
}

func (this *planRepoInMemory) Save(plan *domain.Plan) bool {
	PlansMux.Lock()
	defer PlansMux.Unlock()

	for i := 0; i < len(Plans); i++ {
		if Plans[i].TopicName == plan.TopicName && Plans[i].Title == plan.Title {
			return false
		}
	}
	PlansCounter++
	plan.Id = PlansCounter
	Plans = append(Plans, *plan)
	return true
}

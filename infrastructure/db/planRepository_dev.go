//
//
package db

//
//import (
//	"fmt"
//	"sort"
//	"sync"
//
//	"github.com/NeekUP/roadmaps/core"
//	"github.com/jackc/pgx/v4"
//
//	"github.com/NeekUP/roadmaps/domain"
//)
//
//var (
//	Plans        = make([]domain.Plan, 0)
//	PlansMux     sync.Mutex
//	PlansCounter int
//)
//
//type planRepoInMemory struct {
//	Conn *pgx.Conn
//}
//
//func NewPlansRepository(conn *pgx.Conn) core.PlanRepository {
//	return &planRepoInMemory{
//		Conn: conn}
//}
//
//func (this *planRepoInMemory) SaveWithSteps(plan *domain.Plan) error {
//	TopicsMux.Lock()
//	defer TopicsMux.Unlock()
//
//	var topic *domain.TopicName
//	for i := 0; i < len(Topics); i++ {
//		if Topics[i].Name == plan.TopicName {
//			topic = &Topics[i]
//		}
//	}
//
//	if topic == nil {
//		return fmt.Errorf("TopicName not found be name: %s", plan.TopicName)
//	}
//
//	if !this.Save(plan) {
//		return fmt.Errorf("Plan with same title already exists: %s", plan.Title)
//	}
//
//	StepsMux.Lock()
//	defer StepsMux.Unlock()
//
//	for i := 0; i < len(plan.Steps); i++ {
//		StepsCounter++
//		plan.Steps[i].PlanId = plan.Id
//		plan.Steps[i].Id = StepsCounter
//		Steps = append(Steps, plan.Steps[i])
//	}
//
//	return nil
//}
//
//func (this *planRepoInMemory) Get(id int) *domain.Plan {
//	PlansMux.Lock()
//	defer PlansMux.Unlock()
//
//	for _, v := range Plans {
//		if v.Id == id {
//			// for _, s := range Steps {
//			// 	if s.PlanId == id {
//			// 		result.Steps = append(result.Steps, s)
//			// 	}
//			// }
//
//			copy := v
//			this.AttacheSteps(&copy)
//			return &copy
//		}
//	}
//
//	return nil
//}
//
//func (this *planRepoInMemory) GetList(id []int) []domain.Plan {
//	list := make([]domain.Plan, 0)
//	for i := 0; i < len(id); i++ {
//		p := this.Get(id[i])
//		if p != nil {
//			copy := *p
//			list = append(list, copy)
//		}
//	}
//
//	return list
//}
//
//func (this *planRepoInMemory) GetPopularByTopic(topic string, count int) []domain.Plan {
//
//	PlansMux.Lock()
//	defer PlansMux.Unlock()
//
//	topicPlans := make([]domain.Plan, 0)
//
//	for i := 0; i < PlansCounter; i++ {
//		if Plans[i].TopicName == topic {
//			topicPlans = append(topicPlans, Plans[i])
//		}
//	}
//
//	sort.Slice(topicPlans, func(i, j int) bool {
//		return topicPlans[i].Points > topicPlans[j].Points
//	})
//
//	if count > len(topicPlans) {
//		count = len(topicPlans)
//	}
//
//	result := topicPlans[:count]
//
//	return result
//}
//
//func (this *planRepoInMemory) AttacheSteps(plans *domain.Plan) {
//	StepsMux.Lock()
//	defer StepsMux.Unlock()
//
//	for s := 0; s < len(Steps); s++ {
//		if Steps[s].PlanId != plans.Id {
//			continue
//		}
//
//		step := Steps[s]
//		if step.ReferenceType == domain.ResourceReference {
//			for r := 0; r < len(Sources); r++ {
//				if Sources[r].Id == step.Id {
//					sourceCopy := Sources[r]
//					step.Source = &sourceCopy
//					break
//				}
//			}
//		} else if step.ReferenceType == domain.TopicReference {
//			for t := 0; t < len(Topics); t++ {
//				if step.ReferenceId == Topics[t].Id {
//					topicCopy := Topics[t]
//					step.Source = &domain.Source{
//						Id:         -1,
//						Title:      topicCopy.Title,
//						Identifier: topicCopy.Name,
//						Desc:       topicCopy.Description,
//					}
//				}
//			}
//		}
//		plans.Steps = append(plans.Steps, step)
//	}
//}
//
//func (this *planRepoInMemory) Save(plan *domain.Plan) bool {
//	PlansMux.Lock()
//	defer PlansMux.Unlock()
//
//	for i := 0; i < len(Plans); i++ {
//		if Plans[i].TopicName == plan.TopicName && Plans[i].Title == plan.Title {
//			return false
//		}
//	}
//	PlansCounter++
//	plan.Id = PlansCounter
//	Plans = append(Plans, *plan)
//	return true
//}
//
//func (this *planRepoInMemory) All() []domain.Plan {
//	return Plans
//}

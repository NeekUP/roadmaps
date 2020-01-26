package domain

type EntityType int

const (
	PlanEntity     EntityType = 1
	TopicEntity    EntityType = 2
	ProjectEntity  EntityType = 3
	ResourceEntity EntityType = 4
)

func (et EntityType) IsValid() bool {
	return et >= 1 && et <= 5
}

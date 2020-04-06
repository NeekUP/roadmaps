package domain

import "strings"

type EntityType int

const (
	PlanEntity     EntityType = 1
	TopicEntity    EntityType = 2
	ProjectEntity  EntityType = 3
	ResourceEntity EntityType = 4
	CommentEntity  EntityType = 5
	UserEntity     EntityType = 6
)

func (et EntityType) IsValid() bool {
	return et >= 1 && et <= 6
}

func EntityTypeFromString(entityType string) (bool, EntityType) {
	switch strings.ToLower(entityType) {
	case "plan":
		return true, PlanEntity
	case "topic":
		return true, PlanEntity
	case "project":
		return true, PlanEntity
	case "resource":
		return true, PlanEntity
	case "comment":
		return true, PlanEntity
	case "user":
		return true, PlanEntity
	default:
		return false, 0
	}
}

func EntityTypeToString(entityType EntityType) string {
	switch entityType {
	case PlanEntity:
		return "plan"
	case CommentEntity:
		return "comment"
	case ProjectEntity:
		return "project"
	case TopicEntity:
		return "topic"
	case ResourceEntity:
		return "resource"
	case UserEntity:
		return "user"
	default:
		return ""
	}
}

package domain

type ReferenceType string

const (
	ResourceReference ReferenceType = "Resource"
	TopicReference    ReferenceType = "Topic"
	ProjectReference  ReferenceType = "Project"
)

package domain

type ReferenceType string

const (
	ResourceReference ReferenceType = "Resource"
	TopicReference    ReferenceType = "Topic"
	ProjectReference  ReferenceType = "Project"
)

func (r ReferenceType) IsValid() bool {
	return r == ResourceReference ||
		r == TopicReference ||
		r == ProjectReference
}

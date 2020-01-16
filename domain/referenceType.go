package domain

type ReferenceType string

const (
	ResourceReference ReferenceType = "Resource"
	TopicReference    ReferenceType = "Topic"
	TestReference     ReferenceType = "Test"
)

package domain

type ReferenceType string

const (
	ResourceReference ReferenceType = "Resource"
	TopicReference    ReferenceType = "TopicName"
	TestReference     ReferenceType = "Test"
)

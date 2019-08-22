package domain

type Step struct {
	Id            int
	TopicId       string
	ReferenceId   string
	ReferenceType ReferenceType
	Position      int
}

package domain

type Step struct {
	Id            string
	TopicId       string
	ReferenceId   string
	ReferenceType string // ReferenceType string
	Position      int
}

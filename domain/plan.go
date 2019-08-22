package domain

type Plan struct {
	Id        int
	TopicName string
	OwnerId   string
	Steps     []Step
}

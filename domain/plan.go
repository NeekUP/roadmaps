package domain

type Plan struct {
	Id        int
	Title     string
	TopicName string
	OwnerId   string
	Points    int
	Steps     []Step
}

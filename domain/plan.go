package domain

type Plan struct {
	Id        int
	Title     string
	TopicName string
	OwnerId   string
	Steps     []Step
	Owner     *User
	Points    *Points
	IsDraft   bool
}

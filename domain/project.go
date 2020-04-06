package domain

type Project struct {
	Id      int
	Title   string
	Text    string
	Tags    []TopicTag
	OwnerId string
	Owner   *User
	Points  *Points
}

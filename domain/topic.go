package domain

import "github.com/gosimple/slug"

type Topic struct {
	Id          int
	Name        string
	Title       string
	Description string
	Creator     string
	Plans       []Plan
}

func NewTopic(title, desc, userID string) *Topic {
	return &Topic{
		Name:        generateTopicID(title),
		Title:       title,
		Description: desc,
		Creator:     userID,
	}
}

func generateTopicID(title string) string {
	m := make(map[string]string)
	m["#"] = "-sharp"
	m["+"] = "-p"

	return slug.Make(slug.Substitute(title, m))
}

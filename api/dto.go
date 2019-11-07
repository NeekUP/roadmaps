package api

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type topic struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"desc,omitempty"`
	Plans       []plan `json:"plans,omitempty"`
}

func NewTopicDto(t *domain.Topic) *topic {
	nt := &topic{
		Name:        t.Name,
		Title:       t.Title,
		Description: t.Description,
		Plans:       make([]plan, len(t.Plans)),
	}

	for i := 0; i < len(t.Plans); i++ {
		nt.Plans[i] = *NewPlanDto(&t.Plans[i], i == 0)
	}

	return nt
}

type plan struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	TopicName   string `json:"topicName"`
	Owner       *user  `json:"owner,omitempty"`
	Points      int    `json:"points"`
	InFavorites bool   `json:"inFavorites"`
	Steps       []step `json:"steps,omitempty"`
}

func NewPlanDto(p *domain.Plan, inFavorites bool) *plan {
	np := &plan{
		Id:          core.EncodeNumToString(p.Id),
		Title:       p.Title,
		Points:      p.Points,
		TopicName:   p.TopicName,
		Owner:       NewUserDto(p.Owner),
		InFavorites: inFavorites,
		Steps:       make([]step, len(p.Steps)),
	}

	for i := 0; i < len(p.Steps); i++ {
		np.Steps[i] = *NewStepDto(&p.Steps[i])
	}

	return np
}

type step struct {
	Id            int                  `json:"id"`
	ReferenceType domain.ReferenceType `json:"type"`
	Position      int                  `json:"position"`
	Source        *source              `json:"source"`
}

func NewStepDto(s *domain.Step) *step {
	if s == nil {
		return nil
	}

	return &step{
		Id:            s.Id,
		ReferenceType: s.ReferenceType,
		Position:      s.Position,
		Source:        NewSourceDto(s.Source),
	}
}

type source struct {
	Id         interface{}       `json:"id"`
	Title      string            `json:"title"`
	Type       domain.SourceType `json:"type,omitempty"`
	Properties string            `json:"props,omitempty"`
	Img        string            `json:"img,omitempty"`
	Desc       string            `json:"desc,omitempty"`
}

func NewSourceDto(s *domain.Source) *source {
	if s == nil {
		return nil
	}

	result := &source{
		Id:         s.Id,
		Title:      s.Title,
		Type:       s.Type,
		Properties: s.Properties,
		Img:        s.Img,
		Desc:       s.Desc,
	}

	if s.Id == -1 {
		result.Id = s.Identifier
	}

	return result
}

type user struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

func NewUserDto(u *domain.User) *user {
	if u == nil {
		return nil
	}

	return &user{
		Id:   u.Id,
		Name: u.Name,
		Img:  u.Img,
	}
}
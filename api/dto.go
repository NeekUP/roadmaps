package api

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type topic struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description string     `json:"desc,omitempty"`
	Tags        []topicTag `json:"tags"`
	Plans       []plan     `json:"plans,omitempty"`
	IsTag       bool       `json:"isTag"`
}

func NewTopicDto(t *domain.Topic) *topic {
	nt := &topic{
		Id:          t.Id,
		Name:        t.Name,
		Title:       t.Title,
		Description: t.Description,
		Tags:        make([]topicTag, len(t.Tags)),
		Plans:       make([]plan, len(t.Plans)),
		IsTag:       t.IsTag,
	}

	for i := 0; i < len(t.Tags); i++ {
		nt.Tags[i] = *NewTopicTag(&t.Tags[i])
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
	Id            int64                `json:"id"`
	ReferenceType domain.ReferenceType `json:"type"`
	Position      int                  `json:"position"`
	Source        interface{}          `json:"source"`
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

func NewSourceDto(s interface{}) interface{} {
	if s == nil {
		return nil
	}

	switch v := s.(type) {
	case *domain.Source:
		src := &source{
			Id:         v.Id,
			Title:      v.Title,
			Type:       v.Type,
			Properties: v.Properties,
			Img:        v.Img,
			Desc:       v.Desc,
		}

		if v.Id == -1 {
			src.Id = v.Identifier
		}
		return src
	case *domain.Topic:
		tpc := &topic{
			Id:          v.Id,
			Title:       v.Title,
			Name:        v.Name,
			Description: v.Description,
			IsTag:       v.IsTag,
		}
		return tpc
	}

	return nil
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

type topicTag struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func NewTopicTag(t *domain.TopicTag) *topicTag {
	return &topicTag{
		Name:  t.Name,
		Title: t.Title,
	}
}

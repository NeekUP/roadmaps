package tests

import (
	"roadmaps/core/usecases"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"testing"

	"github.com/google/uuid"
)

func TestAddTopicSuccess(t *testing.T) {

	values := []struct {
		Title, Desc, UserId string
	}{
		{"Параметризация нейросетью физической модели для решения задачи топологической оптимизации", "Параметризация нейросетью физической модели для решения задачи топологической оптимизации", uuid.New().String()},
		{"HTML and CSS: Design and Build Websites", "A full-color introduction to the basics of HTML and CSS from the publishers of Wrox! Every day, more and more people want to learn some HTML and CSS. Joining the professional web designers and programmers are new audiences who need to know a little bit of code at work (update a content management system or e-commerce store) and those who want to make their personal blogs more attractive. Many books teaching HTML and CSS are dry and only written for those who want to become programmers, which is why this", uuid.New().String()},
	}

	usecase := usecases.NewAddTopic(db.NewTopicRepository(nil), log)

	for _, v := range values {
		result, err := usecase.Do(infrastructure.NewContext(nil), v.Title, v.Desc)
		if err != nil {
			t.Errorf("Request with title [%s] return err: %s", v.Title, err.Error())
		}

		if result.Id == 0 {
			t.Errorf("Request with title [%s] hasn't Id", v.Title)
		}

		if result.Description != v.Desc {
			t.Errorf("Request with title [%s] have invalid description: %s ", v.Title, v.Desc)
		}
	}

}

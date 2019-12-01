package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
	"strconv"
)

type addPlanRequest struct {
	TopicName    string        `json:"topic"`
	Title        string        `json:"title"`
	AddPlanSteps []addPlanstep `json:"steps"`
}

type addPlanstep struct {
	ReferenceId   int64                `json:"referenceId"`
	ReferenceType domain.ReferenceType `json:"referenceType"`
}

type addPlanResponse struct {
	TopicName string `json:"topic"`
	Title     string `json:"title"`
	Id        string `json:"id"`
}

func AddPlan(addPlan usecases.AddPlan, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(addPlanRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		addPlanReq := usecases.AddPlanReq{
			Title:     data.Title,
			TopicName: data.TopicName,
		}

		for _, v := range data.AddPlanSteps {
			addPlanReq.Steps = append(addPlanReq.Steps, usecases.PlanStep{ReferenceId: v.ReferenceId, ReferenceType: v.ReferenceType})
		}

		plan, err := addPlan.Do(infrastructure.NewContext(r.Context()), addPlanReq)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &addPlanResponse{
			TopicName: plan.TopicName,
			Title:     plan.Title,
			Id:        core.EncodeNumToString(plan.Id),
		})

	}
}

type getPlanRequest struct {
	Id string `json:"id"`
}

type getPlanTreeResponse struct {
	Nodes []treeNode `json:"nodes"`
}

type treeNode struct {
	TopicName  string     `json:"topicName"`
	TopicTitle string     `json:"topicTitle"`
	PlanId     string     `json:"planId"`
	PlanTitle  string     `json:"planTitle"`
	Child      []treeNode `json:"child"`
}

func GetPlanTree(getPlanTree usecases.GetPlanTree, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			log.Errorf("%s", err.Error())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		// TODO: Remove this
		id, err := strconv.Atoi(data.Id)
		if err != nil {
			id, err = core.DecodeStringToNum(data.Id)
			if err != nil {
				errors := make(map[string]string)
				errors["id"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
		}

		trees, err := getPlanTree.Do(infrastructure.NewContext(r.Context()), []int{id})
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		roots := make([]treeNode, len(trees))
		if len(trees) > 0 {
			for i := 0; i < len(trees); i++ {
				newPlanTree(trees[i], &roots[i])
			}
		}

		valueResponse(w, &getPlanTreeResponse{Nodes: roots})
	}
}

func newPlanTree(node usecases.TreeNode, tree *treeNode) {

	tree.TopicTitle = node.TopicTitle
	tree.TopicName = node.TopicName
	tree.PlanTitle = node.PlanTitle
	tree.PlanId = core.EncodeNumToString(node.PlanId)

	if len(node.Child) > 0 {
		childs := make([]treeNode, len(node.Child))
		for i := 0; i < len(node.Child); i++ {
			childs[i] = treeNode{}
			newPlanTree(node.Child[i], &childs[i])
		}
		tree.Child = childs
	}
}

func GetPlan(getPlan usecases.GetPlan, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		// TODO: Remove this
		id, err := strconv.Atoi(data.Id)
		if err != nil {
			id, err = core.DecodeStringToNum(data.Id)
			if err != nil {
				errors := make(map[string]string)
				errors["id"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
		}

		plan, err := getPlan.Do(infrastructure.NewContext(r.Context()), id)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, NewPlanDto(plan, false))
	}
}

type getPlanListRequest struct {
	TopicName string `json:"topicName"`
}

func GetPlanList(getPlanList usecases.GetPlanList, getUsersPlan usecases.GetUsersPlan, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanListRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		list, err := getPlanList.Do(infrastructure.NewContext(r.Context()), data.TopicName, 100)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		usersPlan, err := getUsersPlan.Do(infrastructure.NewContext(r.Context()), data.TopicName)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		pl := make(map[int]bool)
		result := make([]plan, len(list))
		for i := 0; i < len(list); i++ {
			if pl[list[i].Id] == false {
				pl[list[i].Id] = true
			} else {
				continue
			}

			if usersPlan != nil && usersPlan.Id == list[i].Id {
				result[i] = *NewPlanDto(&list[i], true)
			} else {
				result[i] = *NewPlanDto(&list[i], false)
			}
		}
		valueResponse(w, result)
	}
}

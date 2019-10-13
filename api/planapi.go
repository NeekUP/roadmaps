package api

import (
	"encoding/json"
	"net/http"
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/domain"
	"roadmaps/infrastructure"
	"strconv"
)

type addPlanRequest struct {
	TopicName string `json:"topic"`
	Title     string `json:"title"`
	Steps     []step `json:"steps"`
}

type step struct {
	ReferenceId   int                  `json:"referenceId"`
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

		for _, v := range data.Steps {
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

type getPlanTreeRequest struct {
	Id string
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

func GetPlanTree(addPlan usecases.GetPlanTree, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanTreeRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		id, err := strconv.Atoi(data.Id)
		if err != nil {
			id, err = core.DecodeStringToNum(data.Id)
			if err != nil {
				statusResponse(w, &status{Code: http.StatusBadRequest})
				return
			}
		}

		trees, err := addPlan.Do(infrastructure.NewContext(r.Context()), []int{id})
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

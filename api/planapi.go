package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

type addPlanRequest struct {
	TopicName string     `json:"topic"`
	Title     string     `json:"title"`
	IsDraft   bool       `json:"isDraft"`
	Steps     []planstep `json:"steps"`
}

func (req *addPlanRequest) Sanitize() {
	req.TopicName = StrictSanitize(req.TopicName)
	req.Title = StrictSanitize(req.Title)

	for i := 0; i < len(req.Steps); i++ {
		req.Steps[i].Title = SanitizeText(req.Steps[i].Title)
	}
}

type planstep struct {
	Type   domain.ReferenceType `json:"type"`
	Title  string               `json:"title"`
	Source struct {
		Id int64 `json:"id"`
	}
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
		data.Sanitize()
		addPlanReq := usecases.AddPlanReq{
			Title:     data.Title,
			TopicName: data.TopicName,
			IsDraft:   data.IsDraft,
		}

		for _, v := range data.Steps {
			addPlanReq.Steps = append(addPlanReq.Steps, usecases.PlanStep{ReferenceId: v.Source.Id, ReferenceType: v.Type, Title: v.Title})
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

type editPlanRequest struct {
	Id        string     `json:"id"`
	TopicName string     `json:"topic"`
	Title     string     `json:"title"`
	IsDraft   bool       `json:"isDraft"`
	Steps     []planstep `json:"steps"`
}

func (req *editPlanRequest) Sanitize() {
	req.Id = StrictSanitize(req.Id)
	req.Title = StrictSanitize(req.Title)
	req.TopicName = StrictSanitize(req.TopicName)
	for i := 0; i < len(req.Steps); i++ {
		req.Steps[i].Title = SanitizeText(req.Steps[i].Title)
	}
}

func EditPlan(editPlan usecases.EditPlan, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(editPlanRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		id, err := core.DecodeStringToNum(data.Id)
		if err != nil {
			errors := make(map[string]string)
			errors["id"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		addPlanReq := usecases.EditPlanReq{
			Id:        id,
			Title:     data.Title,
			TopicName: data.TopicName,
			IsDraft:   data.IsDraft,
		}

		for _, v := range data.Steps {
			addPlanReq.Steps = append(addPlanReq.Steps, usecases.PlanStep{ReferenceId: v.Source.Id, ReferenceType: v.Type, Title: v.Title})
		}

		_, err = editPlan.Do(infrastructure.NewContext(r.Context()), addPlanReq)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		statusResponse(w, &status{Code: http.StatusOK})
	}
}

type getPlanRequest struct {
	Id string `json:"id"`
}

func (req *getPlanRequest) Sanitize() {
	req.Id = StrictSanitize(req.Id)
}

type getPlanTreeResponse struct {
	Nodes []treeNode `json:"nodes"`
}

type treeNode struct {
	TopicName  string     `json:"topicName"`
	TopicTitle string     `json:"topicTitle"`
	PlanId     string     `json:"planId"`
	PlanTitle  string     `json:"planTitle"`
	Child      []treeNode `json:"children"`
}

func GetPlanTree(getPlanTree usecases.GetPlanTree, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			log.Errorw("Fail to deserialize getPlanRequest", "error", err.Error())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		id, err := core.DecodeStringToNum(data.Id)
		if err != nil {
			errors := make(map[string]string)
			errors["id"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
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
		child := make([]treeNode, len(node.Child))
		for i := 0; i < len(node.Child); i++ {
			child[i] = treeNode{}
			newPlanTree(node.Child[i], &child[i])
		}
		tree.Child = child
	}
}

func GetPlan(getPlan usecases.GetPlan, getUsersPlan usecases.GetUsersPlan, getPoints usecases.GetPoints, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		id, err := core.DecodeStringToNum(data.Id)
		if err != nil {
			errors := make(map[string]string)
			errors["id"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		ctx := infrastructure.NewContext(r.Context())
		plan, err := getPlan.Do(ctx, id)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		usersPlan, err := getUsersPlan.Do(ctx, plan.TopicName)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		points, err := getPoints.Do(ctx, domain.PlanEntity, int64(plan.Id))
		if err != nil {
			log.Errorw("fail to retrieve points for plan",
				"reqid", ctx.ReqId(),
				"error", "see db log")
		}

		plan.Points = points
		isFavorite := false
		if usersPlan != nil {
			isFavorite = plan.Id == usersPlan.Id
		}
		valueResponse(w, NewPlanDto(plan, isFavorite))
	}
}

type getPlanListRequest struct {
	TopicName string `json:"topicName"`
}

func (req *getPlanListRequest) Sanitize() {
	req.TopicName = StrictSanitize(req.TopicName)
}

func GetPlanList(getPlanList usecases.GetPlanList, getUsersPlan usecases.GetUsersPlan, getPointsList usecases.GetPointsList, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanListRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		data.Sanitize()
		ctx := infrastructure.NewContext(r.Context())
		list, err := getPlanList.Do(ctx, data.TopicName, 100)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		usersPlan, err := getUsersPlan.Do(ctx, data.TopicName)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		idList := make([]int64, len(list))
		for i := 0; i < len(list); i++ {
			idList[i] = int64(list[i].Id)
		}

		points, err := getPointsList.Do(ctx, domain.PlanEntity, idList)
		if err != nil {
			log.Errorw("fail to retrieve points for plan",
				"reqid", ctx.ReqId(),
				"error", "see db log")
		} else {
			for i := 0; i < len(list); i++ {
				for j := 0; j < len(list); j++ {
					if int64(list[j].Id) == points[i].Id {
						list[j].Points = &points[i]
						break
					}
				}
			}
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

type getPlanListByUserReq struct {
	Count int `json:"count"`
	Page  int `json:"page"`
}

func GetListByUser(getByUser usecases.GetPlanListByUser, getPointsList usecases.GetPointsList, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getPlanListByUserReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		ctx := infrastructure.NewContext(r.Context())
		list, err := getByUser.Do(ctx, data.Count, data.Page)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		idList := make([]int64, len(list))
		for i := 0; i < len(list); i++ {
			idList[i] = int64(list[i].Id)
		}

		points, err := getPointsList.Do(ctx, domain.PlanEntity, idList)
		if err != nil {
			log.Errorw("fail to retrieve points for plan",
				"reqid", ctx.ReqId(),
				"error", "see db log")
		} else {
			for i := 0; i < len(list); i++ {
				for j := 0; j < len(list); j++ {
					if int64(list[j].Id) == points[i].Id {
						list[j].Points = &points[i]
						break
					}
				}
			}
		}

		pl := make(map[int]bool)
		result := make([]plan, len(list))
		for i := 0; i < len(list); i++ {
			if pl[list[i].Id] == false {
				pl[list[i].Id] = true
			} else {
				continue
			}

			result[i] = *NewPlanDto(&list[i], false)
		}
		valueResponse(w, result)
	}
}

type removePlanReq struct {
	Id string `json:"id"`
}

func (req *removePlanReq) Sanitize() {
	req.Id = StrictSanitize(req.Id)
}

func RemovePlan(removePlan usecases.RemovePlan, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(removePlanReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		id, err := core.DecodeStringToNum(data.Id)
		if err != nil {
			errors := make(map[string]string)
			errors["id"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		_, err = removePlan.Do(infrastructure.NewContext(r.Context()), id)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &removeTopicTagRes{Removed: true})
	}
}

package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

type publicUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

/*
	Register User
******************************************************************/

type regUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type regUserResp struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	NeedConfirmation bool   `json:"confirmation"`
}

func (req *regUserReq) Sanitize() {
	req.Name = StrictSanitize(req.Name)
	req.Email = StrictSanitize(req.Email)
}

func RegUser(regUsr usecases.RegisterUser, log core.AppLogger, captcha Captcha) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if !captcha.Confirm(r) {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		decoder := json.NewDecoder(r.Body)
		data := new(regUserReq)
		err := decoder.Decode(data)
		defer r.Body.Close()
		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		data.Sanitize()
		u, err := regUsr.Do(infrastructure.NewContext(r.Context()), data.Name, data.Email, data.Pass)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &regUserResp{
			Name:             u.Name,
			Email:            u.Email,
			NeedConfirmation: !u.EmailConfirmed,
		})
	}
}

/*
	Login User
******************************************************************/

type loginUserReq struct {
	Email       string `json:"email"`
	Pass        string `json:"pass"`
	Fingerprint string `json:"fp"`
}

func (req *loginUserReq) Sanitize() {
	req.Email = StrictSanitize(req.Email)
}

type loginUserRes struct {
	AToken string      `json:"atoken"`
	RToken string      `json:"rtoken"`
	User   *publicUser `json:"user"`
}

func Login(loginUsr usecases.LoginUser, log core.AppLogger, captcha Captcha) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if !captcha.Confirm(r) {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		decoder := json.NewDecoder(r.Body)
		data := new(loginUserReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		user, aToken, rToken, err := loginUsr.Do(infrastructure.NewContext(r.Context()), data.Email, data.Pass, data.Fingerprint, r.UserAgent())

		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &loginUserRes{
			User: &publicUser{
				Id:   user.Id,
				Name: user.Name,
				Img:  user.Img},
			AToken: aToken,
			RToken: rToken})
	}
}

/*
	Refresh auth token
******************************************************************/

type refreshTokenReq struct {
	AToken      string `json:"atoken"`
	RToken      string `json:"rtoken"`
	Fingerprint string `json:"fp"`
}

func RefreshToken(refreshToken usecases.RefreshToken, log core.AppLogger, captcha Captcha) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if !captcha.Confirm(r) {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		decoder := json.NewDecoder(r.Body)
		data := new(refreshTokenReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		aToken, rToken, err := refreshToken.Do(infrastructure.NewContext(r.Context()), data.AToken, data.RToken, data.Fingerprint, r.UserAgent())

		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &loginUserRes{
			AToken: aToken,
			RToken: rToken})
	}
}

/*
	TODO: Validate email
******************************************************************/

type userPlanReq struct {
	PlanId string `json:"planId"`
}

func (req *userPlanReq) Sanitize() {
	req.PlanId = StrictSanitize(req.PlanId)
}

type userPlanRes struct {
	Success bool `json:"success"`
}

func AddUserPlan(addUserPlan usecases.AddUserPlan, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(userPlanReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		ctx := infrastructure.NewContext(r.Context())
		planId, err := core.DecodeStringToNum(data.PlanId)
		if err != nil {
			log.Errorw("Bad request", "UserId", ctx.UserId(), "Error", err.Error())
			statusResponse(w, &status{Code: 400, Message: core.InvalidRequest.String()})
			return
		}

		success, err := addUserPlan.Do(ctx, planId)

		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &userPlanRes{Success: success})
	}
}

func RemoveUserPlan(removeUserPlan usecases.RemoveUserPlan, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(userPlanReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		ctx := infrastructure.NewContext(r.Context())
		planId, err := core.DecodeStringToNum(data.PlanId)
		if err != nil {
			log.Errorw("Bad request", "UserId", ctx.UserId(), "Error", err.Error())
			statusResponse(w, &status{Code: 400, Message: core.InvalidRequest.String()})
			return
		}

		success, err := removeUserPlan.Do(ctx, planId)

		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}
		valueResponse(w, &userPlanRes{Success: success})
	}
}

func EmailConfirmation(confirmation usecases.EmailConfirmation, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.URL.Query()["u"]
		if !ok || len(id[0]) < 1 {
			log.Errorw("Parameters with user id is missing", "url", r.URL.String())
			statusResponse(w, &status{Code: 403})
			return
		}

		secret, ok := r.URL.Query()["s"]
		if !ok || len(secret[0]) < 1 {
			log.Errorw("Parameters with secret is missing", "url", r.URL.String())
			statusResponse(w, &status{Code: 403})
			return
		}

		ctx := infrastructure.NewContext(r.Context())
		_, err := confirmation.Do(ctx, id[0], secret[0])
		if err != nil {
			log.Errorw("Bad request", "error", err.Error(), "url", r.URL.String())
			statusResponse(w, &status{Code: 403})
			return
		}

		http.Redirect(w, r, "/login", 302)
	}
}

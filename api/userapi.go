package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

/*
	Register User
******************************************************************/

type regUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type regUserResp struct {
	User             *user `json:"name"`
	NeedConfirmation bool  `json:"confirmation"`
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
		ctx := infrastructure.NewContext(r.Context())
		u, err := regUsr.Do(ctx, data.Name, data.Email, data.Pass)
		if err != nil {
			log.Errorw("usecase err", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &regUserResp{
			User:             NewUserDto(u),
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
	AToken string `json:"atoken"`
	RToken string `json:"rtoken"`
	User   *user  `json:"user"`
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

		ctx := infrastructure.NewContext(r.Context())
		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		user, aToken, rToken, err := loginUsr.Do(ctx, data.Email, data.Pass, data.Fingerprint, r.UserAgent())

		if err != nil {
			log.Errorw("login", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &loginUserRes{
			User:   NewUserDto(user),
			AToken: aToken,
			RToken: rToken})
	}
}

type registerOauthLinkRequest struct {
	ProviderName string `json:"provider"`
	Name         string `json:"name"`
}

type registerOauthLinkResponse struct {
	Url string `json:"url"`
}

func RegisterOAuthLink(checkUser usecases.CheckUser, openAuth core.OpenAuthenticator, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(registerOauthLinkRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		ctx := infrastructure.NewContext(r.Context())

		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		if !openAuth.HasProvider(data.ProviderName) {
			errors := make(map[string]string)
			errors["provider"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		nameExists, err := checkUser.Do(ctx, data.Name)
		if err != nil {
			log.Errorw("checkUser", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusInternalServerError})
			return
		}

		if nameExists {
			errors := make(map[string]string)
			errors["name"] = core.AlreadyExists.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		url, err := openAuth.RegisterLink(data.Name, data.ProviderName)
		if err != nil {
			log.Errorw("oauth registerLink", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &registerOauthLinkResponse{Url: url})
	}
}

type registerOauthRequest struct {
	ProviderName string `json:"provider"`
	Token        string `json:"token"`
	State        string `json:"state"`
	Fingerprint  string `json:"fp"`
}

type registerOauthResp struct {
	User             *user  `json:"name"`
	NeedConfirmation bool   `json:"confirmation"`
	AToken           string `json:"atoken"`
	RToken           string `json:"rtoken"`
}

func RegisterOAuth(reqOauth usecases.RegisterUserOauth, loginOauth usecases.LoginUserOauth, openAuth core.OpenAuthenticator, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(registerOauthRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		ctx := infrastructure.NewContext(r.Context())

		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		if !openAuth.HasProvider(data.ProviderName) {
			errors := make(map[string]string)
			errors["provider"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		name, email, openid, err := openAuth.Auth(data.ProviderName, data.State, data.Token)
		if err != nil {
			log.Errorw("oauth auth", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		user, err := reqOauth.Do(ctx, name, email, data.ProviderName, openid)
		if err != nil {
			log.Errorw("register oauth", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		_, aToken, rToken, err := loginOauth.Do(ctx, data.ProviderName, openid, data.Fingerprint, r.UserAgent())
		if err != nil {
			log.Errorw("login oauth", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
		}
		valueResponse(w, &registerOauthResp{
			User:             NewUserDto(user),
			NeedConfirmation: true,
			AToken:           aToken,
			RToken:           rToken,
		})
	}
}

type loginOauthLinkRequest struct {
	ProviderName string `json:"provider"`
}

type loginOauthLinkResponse struct {
	Url string `json:"url"`
}

func LoginOAuthLink(openAuth core.OpenAuthenticator, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(loginOauthLinkRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()
		ctx := infrastructure.NewContext(r.Context())

		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		if !openAuth.HasProvider(data.ProviderName) {
			errors := make(map[string]string)
			errors["provider"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		url, err := openAuth.LoginLink(data.ProviderName)
		if err != nil {
			log.Errorw("oauth login link", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &loginOauthLinkResponse{Url: url})
	}
}

type loginOauthRequest struct {
	ProviderName string `json:"provider"`
	Token        string `json:"token"`
	State        string `json:"state"`
	Fingerpring  string `json:"fp"`
}

func LoginOauth(loginUsr usecases.LoginUserOauth, openAuth core.OpenAuthenticator, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(loginOauthRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()
		ctx := infrastructure.NewContext(r.Context())

		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		if !openAuth.HasProvider(data.ProviderName) {
			errors := make(map[string]string)
			errors["provider"] = core.InvalidValue.String()
			badRequest(w, core.ValidationError(errors))
			return
		}

		_, _, openid, err := openAuth.Auth(data.ProviderName, data.State, data.Token)
		if err != nil {
			log.Errorw("oauth auth", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		user, aToken, rToken, err := loginUsr.Do(ctx, data.ProviderName, openid, data.Fingerpring, r.UserAgent())
		if err != nil {
			log.Errorw("oauth login", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
		}

		valueResponse(w, &loginUserRes{
			User:   NewUserDto(user),
			AToken: aToken,
			RToken: rToken})
	}
}

type checkUsernameRequest struct {
	Username string `json:"name"`
}

type checkUsernameResponse struct {
	IsFree bool `json:"isFree"`
}

func CheckUser(checkUser usecases.CheckUser, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(checkUsernameRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()
		ctx := infrastructure.NewContext(r.Context())

		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		nameExists, err := checkUser.Do(ctx, data.Username)
		if err != nil {
			log.Errorw("check user", "error", err.Error(), "reqid", ctx.ReqId())
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &checkUsernameResponse{nameExists})
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
		ctx := infrastructure.NewContext(r.Context())

		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		aToken, rToken, err := refreshToken.Do(ctx, data.AToken, data.RToken, data.Fingerprint, r.UserAgent())

		if err != nil {
			log.Errorw("refresh token", "error", err.Error(), "reqid", ctx.ReqId())
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
		ctx := infrastructure.NewContext(r.Context())

		if err != nil {
			log.Errorw("parse request", "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		planId, err := core.DecodeStringToNum(data.PlanId)
		if err != nil {
			log.Errorw("Bad request", "userid", ctx.UserId(), "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: 400, Message: core.InvalidRequest.String()})
			return
		}

		success, err := addUserPlan.Do(ctx, planId)

		if err != nil {
			log.Errorw("add users plan", "error", err.Error(), "reqid", ctx.ReqId())
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
			log.Errorw("Bad request", "userid", ctx.UserId(), "error", err.Error(), "reqid", ctx.ReqId())
			statusResponse(w, &status{Code: 400, Message: core.InvalidRequest.String()})
			return
		}

		success, err := removeUserPlan.Do(ctx, planId)

		if err != nil {
			log.Errorw("remove users plan", "error", err.Error(), "reqid", ctx.ReqId())
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

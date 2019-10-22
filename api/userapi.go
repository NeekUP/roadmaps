package api

import (
	"encoding/json"
	"net/http"
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/infrastructure"
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
		//ua := r.UserAgent()
		_, err = regUsr.Do(infrastructure.NewContext(r.Context()), data.Name, data.Email, data.Pass)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		statusResponse(w, &status{Code: 200})
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

		user, aToken, rToken, err := loginUsr.Do(infrastructure.NewContext(r.Context()), data.Email, data.Pass, data.Fingerprint, r.UserAgent())

		if err != nil {
			if err.Error() != core.InternalError.String() {
				statusResponse(w, &status{Code: 400, Message: err.Error()})
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
				statusResponse(w, &status{Code: 400, Message: err.Error()})
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

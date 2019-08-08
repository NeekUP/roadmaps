package api

import (
	"encoding/json"
	"net/http"
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/infrastructure"
)

type regUserReq struct {
	name  string `json:"name"`
	email string `json:"email"`
	pass  string `json:"pass"`
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
		//TODO: check captcha
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		//ua := r.UserAgent()
		_, err = regUsr.Do(infrastructure.NewContext(r.Context()), data.name, data.email, data.pass)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				statusResponse(w, &status{Code: 400, Message: err.Error()})
			} else {
				statusResponse(w, &status{Code: 500})
			}
		}

	}
}

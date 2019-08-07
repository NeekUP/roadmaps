package api

import "net/http"

type Captcha interface {
	Confirm(r *http.Request) bool
}

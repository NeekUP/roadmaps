package infrastructure

import (
	"net/http"
)

type SuccessCaptcha struct{}

func (SuccessCaptcha) Confirm(r *http.Request) bool {
	return true
}

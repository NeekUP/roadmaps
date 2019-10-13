package tests

import (
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"strings"
	"testing"
)

var (
	pass      string              = "123123"
	fp        string              = "wweqweqweqwwq"
	useragent string              = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"
	ur        core.UserRepository = db.NewUserRepository(nil)
	log       core.AppLogger      = appLoggerForTests{}
	hash      core.HashProvider   = infrastructure.NewSha256HashProvider()
	tokens    core.TokenService   = infrastructure.NewJwtTokenService(ur, "12312312312321")
)

func TestLoginUserSuccess(t *testing.T) {
	email := "TestLoginUserSuccess@123.ww"
	name := "TestLoginUserSuccess"

	regUserAction := usecases.NewRegisterUser(ur, log, hash)
	_, err := regUserAction.Do(infrastructure.NewContext(nil), name, email, pass)
	if err != nil {
		t.Errorf("Register ended with error: %s", err.Error())
		return
	}

	method := usecases.NewLoginUser(ur, log, hash, tokens)

	user, at, rt, err := method.Do(infrastructure.NewContext(nil), email, pass, fp, useragent)
	if user == nil {
		t.Errorf("User is nil")
	}

	if at == "" {
		t.Errorf("Auth token is empty")
	}

	if rt == "" {
		t.Errorf("Refresh token is empty")
	}

	if err != nil {
		t.Errorf("Login ended with error: %s", err.Error())
	}
}

func TestLoginBadPass(t *testing.T) {
	email := "TestLoginBadPass@123.ww"
	name := "TestLoginBadPass"

	regUserAction := usecases.NewRegisterUser(ur, log, hash)
	_, err := regUserAction.Do(infrastructure.NewContext(nil), name, email, pass)
	if err != nil {
		t.Errorf("Register ended with error: %s", err.Error())
		return
	}

	method := usecases.NewLoginUser(ur, log, hash, tokens)

	user, at, rt, err := method.Do(infrastructure.NewContext(nil), email, "3333333", fp, useragent)

	if user != nil {
		t.Errorf("User is not nil")
	}

	if at != "" {
		t.Errorf("Auth token is not empty")
	}

	if rt != "" {
		t.Errorf("Refresh token is not empty")
	}

	if err == nil {
		t.Errorf("Login with bad password ended with no error")
	}

	requestError := strings.Contains(err.Error(), core.AuthenticationError.String())

	if err != nil && !requestError {
		t.Errorf("Unexpected error: %s", err.Error())
	}

}

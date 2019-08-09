package tests

import (
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"testing"
)

func TestRegisterUserSuccess(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(nil), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	email := "TestRegisterUserSuccess@dd.dd"
	name := "TestRegisterUserSuccess"

	user, err := r.Do(infrastructure.FakeContext{}, name, email, "1234")
	if err != nil {
		t.Error(err)
	}

	if user == nil {
		t.Error("user is nil")
		return
	}

	if user.Email != email {
		t.Errorf("not expected user email: [%s]", user.Email)
	}

	if user.Name != name {
		t.Errorf("not expected user name: [%s]", user.Name)
	}
}

func TestRegisterUserInvalidName(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(nil), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	email := "TestRegisterUserInvalidName@dd.dd"
	names := [...]string{
		"name@",
		"",
		"1",
		"asdsd@",
		"email@email.ee",
		"xKb6VMoLWZzYvuihI77kpISTT6QzsS4t1"}

	for i := 0; i < len(names); i++ {
		name := names[i]

		user, err := r.Do(infrastructure.FakeContext{}, name, email, "1234")
		if err == nil {
			t.Error("err is nil")
		}

		if user != nil {
			t.Error("user is not nil")
		}

		if err != nil && err.Error() != core.BadUserName.String() {
			t.Errorf("not expected err: [%s]", err.Error())
		}
	}
}

func TestRegisterUserInvalidEmail(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(nil), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := "TestRegisterUserInvalidEmail"
	emails := [...]string{
		"name@a",
		"name@",
		"",
		"@wwe",
		"d.d@dsdsdsd"}

	for i := 0; i < len(emails); i++ {
		email := emails[i]

		user, err := r.Do(infrastructure.FakeContext{}, name, email, "1234")
		if err == nil {
			t.Errorf("err is nil: [%s]", email)
		}

		if user != nil {
			t.Errorf("user is not nil: [%s]", email)
		}

		if err != nil && err.Error() != core.BadEmail.String() {
			t.Errorf("not expected err: [%s]", err.Error())
		}
	}
}

func TestInvalidPassword(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(nil), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := "TestInvalidPassword"
	email := "TestInvalidPassword@ee.ee"
	passwords := [...]string{
		"name*",
		"***",
		"sd dsd",
		"ssd\"dsd\"dsd",
		"d.d@dsd'sds'd",
		"d.d@dsd|sdsd"}

	for i := 0; i < len(passwords); i++ {
		pass := passwords[i]

		user, err := r.Do(infrastructure.FakeContext{}, name, email, pass)
		if err == nil {
			t.Errorf("err is nil: [%s]", pass)
		}

		if user != nil {
			t.Errorf("user is not nil: [%s]", pass)
		}

		if err != nil && err.Error() != core.BadPassword.String() {
			t.Errorf("not expected err: [%s]", err.Error())
		}
	}
}

func TestRegisterUserExistsName(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(nil), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := "TestRegisterUserExistsName"
	email := "TestRegisterUserExistsName@ee.ee"
	pass := "12345"

	_, err := r.Do(infrastructure.FakeContext{}, name, email, pass)
	if err != nil {
		t.Errorf("Fail to create user")
		return
	}
	user, err := r.Do(infrastructure.FakeContext{}, name, email, pass)

	if user != nil {
		t.Errorf("user is not null")
	}

	if err != nil && err.Error() != core.NameAlreadyExists.String() {
		t.Errorf("not expected err: [%s]", err.Error())
	}
}

func TestRegisterUserExistsEmail(t *testing.T) {
	r := usecases.NewRegisterUser(db.NewUserRepository(nil), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := "TestRegisterUserExistsEmail"
	email := "TestRegisterUserExistsEmail@email.com"
	pass := "12345"

	_, err := r.Do(infrastructure.FakeContext{}, name, email, pass)
	if err != nil {
		t.Errorf("Fail to create user")
		return
	}

	user, err := r.Do(infrastructure.FakeContext{}, "TestRegisterUserExistsEmail_2", email, pass)

	if user != nil {
		t.Errorf("user is not null")
	}

	if err != nil && err.Error() != core.EmailAlreadyExists.String() {
		t.Errorf("not expected err: [%s]", err.Error())
	}
}

///*
//	Context
//*/
//
//type contextForTests struct{}
//
//func (contextForTests) Deadline() (deadline time.Time, ok bool) {
//	panic("implement me")
//}
//
//func (contextForTests) Done() <-chan struct{} {
//	panic("implement me")
//}
//
//func (contextForTests) Err() error {
//	panic("implement me")
//}
//
//func (contextForTests) Value(key interface{}) interface{} {
//	return "123"
//}

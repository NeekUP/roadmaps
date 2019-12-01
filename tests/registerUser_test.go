package tests

import (
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/infrastructure"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"strings"
	"testing"
)

func TestRegisterUserSuccess(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := RandString(10)
	email := fmt.Sprintf("%s@dd.dd", name)
	pass := RandString(10)

	user, err := r.Do(infrastructure.NewContext(nil), name, email, pass)
	defer DeleteUser(user.Id)

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

	r := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

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

		user, err := r.Do(infrastructure.NewContext(nil), name, email, "1234")
		if user != nil {
			defer DeleteUser(user.Id)
		}
		if err == nil {
			t.Error("err is nil")
		}

		if user != nil {
			t.Error("user is not nil")
		}

		formatError := strings.Contains(err.Error(), core.InvalidFormat.String())
		requestError := strings.Contains(err.Error(), core.InvalidRequest.String())

		if err != nil && (!formatError || !requestError) {
			t.Errorf("not expected err: [%s]", err.Error())
		}
	}
}

func TestRegisterUserInvalidEmail(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := "TestRegisterUserInvalidEmail"
	emails := [...]string{
		"name@a",
		"name@",
		"",
		"@wwe",
		"d.d@dsdsdsd"}

	for i := 0; i < len(emails); i++ {
		email := emails[i]

		user, err := r.Do(infrastructure.NewContext(nil), name, email, "1234")
		if user != nil {
			defer DeleteUser(user.Id)
		}
		if err == nil {
			t.Errorf("err is nil: [%s]", email)
		}

		if user != nil {
			t.Errorf("user is not nil: [%s]", email)
		}

		formatError := strings.Contains(err.Error(), core.InvalidFormat.String())
		requestError := strings.Contains(err.Error(), core.InvalidRequest.String())

		if err != nil && (!formatError || !requestError) {
			t.Errorf("not expected err: [%s]", err.Error())
		}
	}
}

func TestInvalidPassword(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

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

		user, err := r.Do(infrastructure.NewContext(nil), name, email, pass)
		if user != nil {
			defer DeleteUser(user.Id)
		}
		if err == nil {
			t.Errorf("err is nil: [%s]", pass)
		}

		if user != nil {
			t.Errorf("user is not nil: [%s]", pass)
		}

		formatError := strings.Contains(err.Error(), core.InvalidFormat.String())
		requestError := strings.Contains(err.Error(), core.InvalidRequest.String())

		if err != nil && (!formatError || !requestError) {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	}
}

func TestRegisterUserExistsName(t *testing.T) {

	r := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := "TestRegisterUserExistsName"
	email := "TestRegisterUserExistsName@ee.ee"
	pass := "12345"

	u, err := r.Do(infrastructure.NewContext(nil), name, email, pass)
	if u != nil {
		defer DeleteUser(u.Id)
	}
	if err != nil {
		t.Errorf("Fail to create user")
		return
	}
	user, err := r.Do(infrastructure.NewContext(nil), name, email, pass)

	if user != nil {
		t.Errorf("user is not null")
	}

	if err != nil && strings.Contains(core.AlreadyExists.String(), err.Error()) {
		t.Errorf("not expected err: [%s]", err.Error())
	}
}

func TestRegisterUserExistsEmail(t *testing.T) {
	r := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})

	name := "TestRegisterUserExistsEmail"
	email := "TestRegisterUserExistsEmail@email.com"
	pass := "12345"

	u, err := r.Do(infrastructure.NewContext(nil), name, email, pass)
	if u != nil {
		defer DeleteUser(u.Id)
	}

	if err != nil {
		t.Errorf("Fail to create user")
		return
	}

	user, err := r.Do(infrastructure.NewContext(nil), "TestRegisterUserExistsEmail_2", email, pass)

	if user != nil {
		t.Errorf("user is not null")
	}

	if err != nil && strings.Contains(core.AlreadyExists.String(), err.Error()) {
		t.Errorf("not expected err: [%s]", err.Error())
	}
}

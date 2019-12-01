package tests

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/infrastructure"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"testing"
)

var (
	secret    string = "ih7Cp1aB0exNXzsHjV9Z66qBczoG8g15_bBBW7iK1L-szDYVIbhWDZv6R-d_PD_TOjriomFr44UYMky2snKInO_7UL23uBmsH6hFlaqGJv12SQl4LC_1D7DW1iNLWSB22u1f3YowVH8YS_odqsUs5klaR7BlsvnQxucJcqSom6JuuZynz3j8p-8MevBDWTPAD7QeD4NUjTp55JftBEEg8J3Qf0ZrFOxkP2ULKvX-VbTwBN2U3YnNHJsdQ5aleUH-62NiG9EUiEDrLuEWw73oHaSCDPLVhIM1zCHW25Nmy8oxzW7rBVPwyLHC9v63QBSH7JXVhBOfDm-F55eOG0zlBw"
	jwtTokens core.TokenService
)

func init() {
	jwtTokens = &infrastructure.JwtTokenService{UserRepo: db.NewUserRepository(DB), Secret: secret}
}

func TestCreateValidateRefreshSuccess(t *testing.T) {
	regUser := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})
	email := RandString(10) + "@test.com"
	user, _ := regUser.Do(infrastructure.NewContext(nil), RandString(10), email, pass)
	a, r, err := jwtTokens.Create(user, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	if a == "" || r == "" || err != nil {
		t.Errorf("Save token return error: [%s]. authToken: [%s] refreshToken: [%s]", err.Error(), a, r)
	}

	if user.Tokens == nil || len(user.Tokens) != 1 {
		t.Errorf("Refresh token id not stored into user info")
	}

	uid, _, rights, err := jwtTokens.Validate(a)

	if err != nil {
		t.Errorf("Error whilw validation token: [%s]", err.Error())
	}

	if rights != 1 {
		t.Errorf("Rights from auth token invalid: [%s].Rights:%d but expected:%d", uid, rights, 1)
	}

	aa, rr, err := jwtTokens.Refresh(a, r, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	if aa == "" || rr == "" || err != nil {
		t.Errorf("refresh token return error: [%s]. authToken: [%s] refreshToken: [%s]", err.Error(), a, r)
	}

	if aa == a || rr == r {
		t.Errorf("Tokens after and before refresh a equals")
	}

	if user.Tokens == nil || len(user.Tokens) != 1 {
		t.Errorf("Refresh token id count not expected [%d]", len(user.Tokens))
	}
}

func TestCreateValidateBadToken(t *testing.T) {
	regUser := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})
	email := RandString(10) + "@test.com"
	user, _ := regUser.Do(infrastructure.NewContext(nil), RandString(10), email, pass)
	a, r, err := jwtTokens.Create(user, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	if a == "" || r == "" || err != nil {
		t.Errorf("Save token return error: [%s]. authToken: [%s] refreshToken: [%s]", err.Error(), a, r)
	}

	bytes := []byte(a)
	bytes[1] = bytes[1] + 1
	a = string(bytes)

	uid, _, _, err := jwtTokens.Validate(a)

	if err == nil {
		t.Errorf("Bad token has been validating")
	}

	if uid != "" {
		t.Errorf("Uid exists but validating fail: [%s]", uid)
	}
}

func TestCreateRefreshByAuthToken(t *testing.T) {
	regUser := usecases.NewRegisterUser(db.NewUserRepository(DB), &appLoggerForTests{}, &infrastructure.Sha256HashProvider{})
	email := RandString(10) + "@test.com"
	user, _ := regUser.Do(infrastructure.NewContext(nil), RandString(10), email, pass)
	a, r, err := jwtTokens.Create(user, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	if a == "" || r == "" || err != nil {
		t.Errorf("Save token return error: [%s]. authToken: [%s] refreshToken: [%s]", err.Error(), a, r)
	}

	aa, rr, err := jwtTokens.Refresh(a, a, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")

	if err == nil {
		t.Errorf("Bad token has been validating")
	}

	if aa != "" {
		t.Errorf("Token returned refresh validating fail: [%s]", aa)
	}

	if rr != "" {
		t.Errorf("Token returned refresh validating fail: [%s]", rr)
	}
}

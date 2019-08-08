package tests

import (
	"roadmaps/core"
	"roadmaps/domain"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"testing"
)

var (
	userId    string = "92382cd6-2acd-431b-b2b4-3500eaa9e342"
	secret    string = "ih7Cp1aB0exNXzsHjV9Z66qBczoG8g15_bBBW7iK1L-szDYVIbhWDZv6R-d_PD_TOjriomFr44UYMky2snKInO_7UL23uBmsH6hFlaqGJv12SQl4LC_1D7DW1iNLWSB22u1f3YowVH8YS_odqsUs5klaR7BlsvnQxucJcqSom6JuuZynz3j8p-8MevBDWTPAD7QeD4NUjTp55JftBEEg8J3Qf0ZrFOxkP2ULKvX-VbTwBN2U3YnNHJsdQ5aleUH-62NiG9EUiEDrLuEWw73oHaSCDPLVhIM1zCHW25Nmy8oxzW7rBVPwyLHC9v63QBSH7JXVhBOfDm-F55eOG0zlBw"
	jwtTokens core.TokenService
	user      *domain.User = &domain.User{Id: userId,
		Name:           "name",
		Email:          "email@email.ru",
		EmailConfirmed: true,
		Img:            "",
		Tokens:         nil,
		Rights:         1,
		Pass:           nil,
		Salt:           nil}
)

func init() {
	userRepo := db.NewUserRepository(nil)
	userRepo.Create(user, nil, nil)
	jwtTokens = &infrastructure.JwtTokenService{UserRepo: userRepo, Secret: secret}
}

func TestCreateValidateRefreshSuccess(t *testing.T) {

	a, r, err := jwtTokens.Create(user, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	if a == "" || r == "" || err != nil {
		t.Errorf("Create token return error: [%s]. authToken: [%s] refreshToken: [%s]", err.Error(), a, r)
	}

	if user.Tokens == nil || len(user.Tokens) != 1 {
		t.Errorf("Refresh token id not stored into user info")
	}

	uid, rights, err := jwtTokens.Validate(a)

	if err != nil {
		t.Errorf("Error whilw validation token: [%s]", err.Error())
	}

	if uid != userId {
		t.Errorf("User Id from auth token invalid: [%s]", uid)
	}

	if rights != 1 {
		t.Errorf("Rights from auth token invalid: [%s]", uid)
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

	a, r, err := jwtTokens.Create(user, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	if a == "" || r == "" || err != nil {
		t.Errorf("Create token return error: [%s]. authToken: [%s] refreshToken: [%s]", err.Error(), a, r)
	}

	bytes := []byte(a)
	bytes[1] = bytes[1] + 1
	a = string(bytes)

	uid, _, err := jwtTokens.Validate(a)

	if err == nil {
		t.Errorf("Bad token has been validating")
	}

	if uid != "" {
		t.Errorf("Uid exists but validating fail: [%s]", uid)
	}
}

func TestCreateRefreshByAuthToken(t *testing.T) {

	a, r, err := jwtTokens.Create(user, "fingerprint", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	if a == "" || r == "" || err != nil {
		t.Errorf("Create token return error: [%s]. authToken: [%s] refreshToken: [%s]", err.Error(), a, r)
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

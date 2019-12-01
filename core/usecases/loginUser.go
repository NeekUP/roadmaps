package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type LoginUser interface {
	Do(tx core.ReqContext, email, password, fingerprint, useragent string) (*domain.User, string, string, error)
}

func NewLoginUser(ur core.UserRepository, log core.AppLogger, hash core.HashProvider, ts core.TokenService) LoginUser {
	return &loginUser{UserRepo: ur, Log: log, Hash: hash, TokenService: ts}
}

type loginUser struct {
	UserRepo     core.UserRepository
	Log          core.AppLogger
	Hash         core.HashProvider
	TokenService core.TokenService
}

func (this *loginUser) Do(ctx core.ReqContext, email, password, fingerprint, useragent string) (*domain.User, string, string, error) {
	appErr := this.validate(ctx, email, password)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"reqId", ctx.ReqId(),
			"email", email,
			"error", appErr.Error(),
		)
		return nil, "", "", appErr
	}

	useragent = core.UserAgentFingerprint(useragent)

	user := this.UserRepo.FindByEmail(email)
	if user == nil {
		this.Log.Infow("User not found",
			"reqId", ctx.ReqId(),
			"email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	if !this.Hash.CheckPassword(password, user.Pass, user.Salt) {
		this.Log.Infow("Password is wrong",
			"reqId", ctx.ReqId(),
			"email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	aToken, rToken, err := this.TokenService.Create(user, fingerprint, useragent)
	if err != nil {
		this.Log.Errorw("Fail to create token pair",
			"reqId", ctx.ReqId(),
			"email", email,
			"error", err.Error())
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	return user, aToken, rToken, nil
}

func (r *loginUser) validate(ctx core.ReqContext, email string, password string) *core.AppError {

	errors := make(map[string]string)
	if !core.IsValidEmail(email) {
		errors["email"] = core.InvalidFormat.String()
	}

	if !core.IsValidPassword(password) {
		errors["pass"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}

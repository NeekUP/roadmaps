package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type LoginUser interface {
	Do(tx core.ReqContext, email, password, fingerprint, useragent string) (*domain.User, string, string, error)
}

func NewLoginUser(ur core.UserRepository, log core.AppLogger, hash core.HashProvider, ts core.TokenService) LoginUser {
	return &loginUser{userRepo: ur, log: log, hash: hash, tokenService: ts}
}

type loginUser struct {
	userRepo     core.UserRepository
	log          core.AppLogger
	hash         core.HashProvider
	tokenService core.TokenService
}

func (usecase *loginUser) Do(ctx core.ReqContext, email, password, fingerprint, useragent string) (*domain.User, string, string, error) {
	appErr := usecase.validate(ctx, email, password)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"reqId", ctx.ReqId(),
			"email", email,
			"error", appErr.Error(),
		)
		return nil, "", "", appErr
	}

	useragent = core.UserAgentFingerprint(useragent)

	user := usecase.userRepo.FindByEmail(email)
	if user == nil {
		usecase.log.Infow("User not found",
			"reqId", ctx.ReqId(),
			"email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	if !usecase.hash.CheckPassword(password, user.Pass, user.Salt) {
		usecase.log.Infow("Password is wrong",
			"reqId", ctx.ReqId(),
			"email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	aToken, rToken, err := usecase.tokenService.Create(user, fingerprint, useragent)
	if err != nil {
		usecase.log.Errorw("Fail to create token pair",
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

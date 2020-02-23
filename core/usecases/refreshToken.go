package usecases

import (
	"github.com/NeekUP/roadmaps/core"
)

type RefreshToken interface {
	Do(ctx core.ReqContext, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error)
}

func NewRefreshToken(ur core.UserRepository, log core.AppLogger, ts core.TokenService, secret string) RefreshToken {
	return &refreshToken{
		userRepo:     ur,
		log:          log,
		secret:       secret,
		tokenService: ts,
	}
}

type refreshToken struct {
	userRepo     core.UserRepository
	log          core.AppLogger
	secret       string
	tokenService core.TokenService
}

func (usecase *refreshToken) Do(ctx core.ReqContext, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error) {
	trace := ctx.StartTrace("refreshToken")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(authToken, refreshToken, fingerprint, useragent)
	if appErr != nil {
		usecase.log.Infow("Not valid data",
			"reqId", ctx.ReqId(),
			"authToken", authToken,
			"refreshToken", refreshToken,
			"fingerprint", fingerprint,
			"useragent", useragent,
			"error", err.Error())
		return "", "", core.NewError(core.InvalidRequest)
	}

	useragent = core.UserAgentFingerprint(useragent)

	auth, refresh, err := usecase.tokenService.Refresh(ctx, authToken, refreshToken, fingerprint, useragent)

	if err != nil {
		usecase.log.Infow("Fail to refresh token",
			"reqId", ctx.ReqId(),
			"authToken", authToken,
			"refreshToken", refreshToken,
			"fingerprint", fingerprint,
			"useragent", useragent,
			"error", err)
		return "", "", core.NewError(core.InternalError)
	}
	return auth, refresh, err
}

func (usecase *refreshToken) validate(aToken, rToken, fingerprint, useragent string) *core.AppError {

	errors := make(map[string]string)
	if !core.IsValidTokenFormat(aToken) {
		errors["atoken"] = core.InvalidFormat.String()
	}

	if !core.IsValidTokenFormat(rToken) {
		errors["rtoken"] = core.InvalidFormat.String()
	}

	if !core.IsValidFingerprint(fingerprint) {
		errors["fp"] = core.InvalidFormat.String()
	}

	if !core.IsValidUserAgent(useragent) {
		errors["useragent"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}

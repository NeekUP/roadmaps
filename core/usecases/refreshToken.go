package usecases

import (
	"roadmaps/core"
)

type RefreshToken interface {
	Do(ctx core.ReqContext, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error)
}

func NewRefreshToken(ur core.UserRepository, log core.AppLogger, ts core.TokenService, secret string) RefreshToken {
	return &refreshToken{
		UserRepo:     ur,
		Log:          log,
		Secret:       secret,
		TokenService: ts,
	}
}

type refreshToken struct {
	UserRepo     core.UserRepository
	Log          core.AppLogger
	Secret       string
	TokenService core.TokenService
}

func (this *refreshToken) Do(ctx core.ReqContext, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error) {

	appErr := this.validate(authToken, refreshToken, fingerprint, useragent)
	if appErr != nil {
		this.Log.Infow("Not valid data",
			"reqId", ctx.ReqId(),
			"authToken", authToken,
			"refreshToken", refreshToken,
			"fingerprint", fingerprint,
			"useragent", useragent,
			"error", err.Error())
		return "", "", core.NewError(core.InvalidRequest)
	}

	useragent = core.UserAgentFingerprint(useragent)

	auth, refresh, err := this.TokenService.Refresh(authToken, refreshToken, fingerprint, useragent)

	if err != nil {
		this.Log.Infow("Fail to refresh token",
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

func (this *refreshToken) validate(aToken, rToken, fingerprint, useragent string) *core.AppError {

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

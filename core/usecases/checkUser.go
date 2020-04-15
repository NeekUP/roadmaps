package usecases

import "github.com/NeekUP/roadmaps/core"

type CheckUser interface {
	Do(ctx core.ReqContext, name string) (bool, error)
}

type checkUser struct {
	usersRepo core.UserRepository
	log       core.AppLogger
}

func NewCheckUser(usersRepo core.UserRepository, log core.AppLogger) CheckUser {
	return &checkUser{
		usersRepo: usersRepo,
		log:       log,
	}
}

func (usecase *checkUser) Do(ctx core.ReqContext, name string) (bool, error) {
	trace := ctx.StartTrace("editProject")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(name)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return false, appErr
	}

	exists, ok := usecase.usersRepo.ExistsName(ctx, name)
	return exists && ok, nil
}

func (usecase *checkUser) validate(name string) *core.AppError {
	errors := make(map[string]string)

	if !core.IsValidUserName(name) {
		errors["name"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type EmailConfirmation interface {
	Do(ctx core.ReqContext, id, secret string) (*domain.User, error)
}

type emailConfirmation struct {
	userRepo core.UserRepository
	log      core.AppLogger
}

func NewEmailConfirmation(userRepo core.UserRepository, log core.AppLogger) EmailConfirmation {
	return &emailConfirmation{
		userRepo: userRepo,
		log:      log,
	}
}

func (usecase *emailConfirmation) Do(ctx core.ReqContext, id, secret string) (*domain.User, error) {
	trace := ctx.StartTrace("emailConfirmation")
	defer ctx.StopTrace(trace)

	user := usecase.userRepo.Get(ctx, id)
	if user == nil {
		return nil, core.NewError(core.AccessDenied)
	}

	if user.EmailConfirmation != secret {
		return nil, core.NewError(core.AccessDenied)
	}

	user.EmailConfirmed = true
	user.EmailConfirmation = ""

	usecase.userRepo.Update(ctx, user)
	return user, nil
}

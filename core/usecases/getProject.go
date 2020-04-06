package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetProject interface {
	Do(ctx core.ReqContext, id int) (*domain.Project, error)
}

type getProject struct {
	projectRepo core.ProjectsRepository
	log         core.AppLogger
}

func NewGetProject(projectRepo core.ProjectsRepository, log core.AppLogger) GetProject {
	return &getProject{
		projectRepo: projectRepo,
		log:         log,
	}
}

func (usecase *getProject) Do(ctx core.ReqContext, id int) (*domain.Project, error) {
	trace := ctx.StartTrace("getProject")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(id)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	p := usecase.projectRepo.Get(ctx, id)
	if p == nil {
		return nil, core.NewError(core.NotExists)
	}

	return p, nil
}

func (usecase *getProject) validate(id int) *core.AppError {
	errors := make(map[string]string)
	if id <= 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

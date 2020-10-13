package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddPermissions interface {
	Do(ctx core.ReqContext, userid string, entityType domain.EntityType, entityid int64, value domain.Permissions) (bool, error)
}

type addPermissions struct {
	repo core.PermissionsRepository
	log  core.AppLogger
}

func (usecase *addPermissions) Do(ctx core.ReqContext, userid string, entityType domain.EntityType, entityid int64, value domain.Permissions) (bool, error) {
	appErr := usecase.validate(entityType, entityid, value)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return false, appErr
	}

	appErr = usecase.repo.Set(ctx, userid, entityType, entityid, value)
	if appErr != nil {
		return false, appErr
	}
	return true, nil
}

func NewAddPermissions(repo core.PermissionsRepository, log core.AppLogger) AddPermissions {
	return &addPermissions{repo: repo, log: log}
}

func (usecase addPermissions) validate(entityType domain.EntityType, entityId int64, value domain.Permissions) *core.AppError {
	errors := make(map[string]string)

	if !core.IsValidEntityType(entityType) {
		errors["type"] = core.InvalidValue.String()
	}

	if entityId <= 0 {
		errors["entityId"] = core.InvalidValue.String()
	}

	if value <= 0 {
		errors["value"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

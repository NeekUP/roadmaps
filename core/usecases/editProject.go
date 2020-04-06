package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type EditProject interface {
	Do(ctx core.ReqContext, id int, title, text string, tags []string) (*domain.Project, error)
}

type editProject struct {
	projectRepo core.ProjectsRepository
	topicRepo   core.TopicRepository
	log         core.AppLogger
	changeLog   core.ChangeLog
}

func NewEditProject(projectRepo core.ProjectsRepository, topicRepo core.TopicRepository, changeLog core.ChangeLog, log core.AppLogger) EditProject {
	return &editProject{projectRepo: projectRepo, topicRepo: topicRepo, changeLog: changeLog, log: log}
}

func (usecase *editProject) Do(ctx core.ReqContext, id int, title, text string, tags []string) (*domain.Project, error) {
	trace := ctx.StartTrace("editProject")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(title, text, tags)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	projectBefore := usecase.projectRepo.Get(ctx, id)
	if projectBefore == nil {
		return nil, core.NewError(core.NotExists)
	}
	userId := ctx.UserId()
	if userId != projectBefore.OwnerId {
		return nil, core.NewError(core.AccessDenied)
	}

	projectAfter := *projectBefore
	projectAfter.Tags = usecase.topicRepo.GetTags(ctx, tags)
	_, err := usecase.projectRepo.Update(ctx, &projectAfter)

	if err != nil {
		usecase.log.Errorw("fail to update project",
			"reqid", ctx.ReqId(),
			"error", err.Error(),
		)
		return nil, err
	}
	usecase.changeLog.Edited(domain.ProjectEntity, int64(projectBefore.Id), userId, projectBefore, &projectAfter)
	return &projectAfter, nil
}

func (usecase *editProject) validate(title, text string, tags []string) *core.AppError {
	errors := make(map[string]string)

	if !core.IsValidProjectTitle(title) {
		errors["title"] = core.InvalidFormat.String()
	}

	if !core.IsValidProjectText(text) {
		errors["text"] = core.InvalidFormat.String()
	}

	for _, tag := range tags {
		if !core.IsValidTopicName(tag) {
			errors["tags"] = core.InvalidFormat.String()
			break
		}
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

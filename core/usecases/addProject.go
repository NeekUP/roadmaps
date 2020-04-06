package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddProject interface {
	Do(ctx core.ReqContext, title, text string, tags []string) (*domain.Project, error)
}

type addProject struct {
	projectRepo core.ProjectsRepository
	topicRepo   core.TopicRepository
	log         core.AppLogger
	changeLog   core.ChangeLog
}

func NewAddProject(projectRepo core.ProjectsRepository, topicRepo core.TopicRepository, changeLog core.ChangeLog, log core.AppLogger) AddProject {
	return &addProject{projectRepo: projectRepo, topicRepo: topicRepo, changeLog: changeLog, log: log}
}

func (usecase *addProject) Do(ctx core.ReqContext, title, text string, tags []string) (*domain.Project, error) {
	trace := ctx.StartTrace("addProject")
	defer ctx.StopTrace(trace)
	appErr := usecase.validate(title, text, tags)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	project := &domain.Project{
		Title:   title,
		Text:    text,
		OwnerId: ctx.UserId(),
	}
	project.Tags = usecase.topicRepo.GetTags(ctx, tags)
	_, err := usecase.projectRepo.Add(ctx, project)

	if err != nil {
		usecase.log.Errorw("fail to save project",
			"reqid", ctx.ReqId(),
			"error", err.Error(),
		)
		return nil, err
	}

	usecase.changeLog.Added(domain.ProjectEntity, int64(project.Id), ctx.UserId())
	return project, nil
}

func (usecase *addProject) validate(title, text string, tags []string) *core.AppError {
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

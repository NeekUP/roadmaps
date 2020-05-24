package db

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type projectsRepo struct {
	Db *DbConnection
}

func NewProjectsRepository(db *DbConnection) core.ProjectsRepository {
	return &projectsRepo{
		Db: db}
}

func (repo *projectsRepo) Add(ctx core.ReqContext, project *domain.Project) (bool, error) {
	panic("implement me")
}

func (repo *projectsRepo) Update(ctx core.ReqContext, project *domain.Project) (bool, error) {
	panic("implement me")
}

func (repo *projectsRepo) Get(ctx core.ReqContext, id int) *domain.Project {
	panic("implement me")
}

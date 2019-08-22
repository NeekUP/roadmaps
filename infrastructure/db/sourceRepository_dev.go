// +build DEV

package db

import (
	"database/sql"
	"roadmaps/core"
	"roadmaps/domain"
	"sync"
)

var (
	Sources        = make([]domain.Source, 0)
	SourcesMux     sync.Mutex
	SourcesCounter int
)

type sourceRepoInMemory struct {
	Conn *sql.DB
}

func NewSourceRepository(conn *sql.DB) core.SourceRepository {
	return &sourceRepoInMemory{
		Conn: conn}
}

func (this *sourceRepoInMemory) Get(id int) {
	panic("not implemented")
}

func (this *sourceRepoInMemory) FindByIdentifier(identifier string) {
	panic("not implemented")
}

func (this *sourceRepoInMemory) Save(source *domain.Source) {
	panic("not implemented")
}

func (this *sourceRepoInMemory) Update(source *domain.Source) {
	panic("not implemented")
}

func (this *sourceRepoInMemory) GetOrAddByIdentifier(source *domain.Source) *domain.Source {
	SourcesMux.Lock()
	defer SourcesMux.Unlock()

	for i := 0; i < len(Sources); i++ {
		if Sources[i].NormalizedIdentifier == source.NormalizedIdentifier {
			return &Sources[i]
		}
	}

	SourcesCounter++
	source.Id = SourcesCounter

	Sources = append(Sources, *source)
	return source
}

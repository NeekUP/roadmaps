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

func (this *sourceRepoInMemory) Get(id int) *domain.Source {
	SourcesMux.Lock()
	defer SourcesMux.Unlock()

	for i := 0; i < len(Users); i++ {
		if Sources[i].Id == id {
			return &Sources[i]
		}
	}
	return nil
}

func (this *sourceRepoInMemory) FindByIdentifier(nIdentifier string) *domain.Source {
	SourcesMux.Lock()
	defer SourcesMux.Unlock()

	for i := 0; i < len(Sources); i++ {
		if Sources[i].NormalizedIdentifier == nIdentifier {
			return &Sources[i]
		}
	}
	return nil
}

func (this *sourceRepoInMemory) Save(source *domain.Source) bool {
	SourcesMux.Lock()
	defer SourcesMux.Unlock()

	for i := 0; i < len(Sources); i++ {
		if Sources[i].NormalizedIdentifier == source.NormalizedIdentifier {
			return false
		}
	}

	Sources = append(Sources, *source)
	return true
}

func (this *sourceRepoInMemory) Update(source *domain.Source) bool {
	SourcesMux.Lock()
	defer SourcesMux.Unlock()

	for i := 0; i < len(Sources); i++ {
		if Sources[i].Id == source.Id {
			Sources[i] = *source
			return true
		}
	}

	return false
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

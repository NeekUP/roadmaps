package usecases

import (
	"fmt"
	"net/url"
	"roadmaps/core"
	"roadmaps/domain"
	"strings"

	"github.com/moraes/isbn"
)

type AddSource interface {
	Do(ctx core.ReqContext, identifier, title, props string, sourceType domain.SourceType) (*domain.Source, error)
}

func NewAddSource(sr core.SourceRepository, log core.AppLogger) AddSource {
	return &addSource{Repo: sr, Log: log}
}

type addSource struct {
	Repo core.SourceRepository
	Log  core.AppLogger
}

// Должен возвращать или уже созданный ранее или новый объект
func (this *addSource) Do(ctx core.ReqContext, identifier, title, props string, sourceType domain.SourceType) (*domain.Source, error) {
	err := this.validate(identifier, title, props, sourceType)
	if err != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return nil, err
	}
	var nIdentifier string

	switch sourceType {
	case domain.Book:
		nIdentifier = this.getBookIdentifier(identifier)
		// TODO: get preview
		break
	case domain.Article:
		nIdentifier, _ = this.getLinkIdentifier(identifier)
		// TODO: get preview
		break
	case domain.Video:
		nIdentifier, _ = this.getLinkIdentifier(identifier)
		// TODO: get preview
		break
	case domain.Audio:
		nIdentifier, _ = this.getLinkIdentifier(identifier)
		// TODO: get preview
		break
	default:
		return nil, core.NewError(core.InvalidSourceType)
	}

	s := &domain.Source{
		Title:                title,
		Identifier:           identifier,
		NormalizedIdentifier: nIdentifier,
		Type:                 sourceType,
		Properties:           "{}"} // TODO: set props

	s = this.Repo.GetOrAddByIdentifier(s)
	return s, nil
}

func (this *addSource) getLinkIdentifier(identifier string) (string, error) {

	identifier = strings.ToLower(identifier)
	identifier = strings.TrimRight(identifier, "/")

	u, err := url.Parse(identifier)
	if err != nil {
		return "", err
	}
	query, err := url.QueryUnescape(u.RawQuery)
	if err != nil {
		return "", err
	}

	if len(query) > 0 {
		query = fmt.Sprintf("?%s", query)
	}

	path, err := url.PathUnescape(u.Path)
	if err != nil {
		return "", err
	}

	host := strings.TrimLeft(u.Host, "www.")
	return fmt.Sprintf("%s%s%s", host, path, query), nil
}

func (this *addSource) getBookIdentifier(identifier string) string {
	b := strings.Replace(identifier, "-", "", -1)
	if len(b) == 10 {
		identifier, _ = isbn.To13(b)
	} else {
		identifier = b
	}
	return identifier
}

func (this *addSource) validate(identifier, title, props string, sourceType domain.SourceType) error {

	if sourceType != domain.Book {
		u, err := url.Parse(identifier)
		if err != nil {
			return err
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return core.NewError(core.InvalidUrl)
		}

		if u.Host == "" {
			return core.NewError(core.InvalidUrl)
		}
	} else {
		b := strings.Replace(identifier, "-", "", -1)
		if !isbn.Validate(b) {
			return core.NewError(core.InvalidISBN)
		}
	}

	if ok, _ := core.IsValidSourceTitle(title); !ok {
		return core.NewError(core.InvalidTitle)
	}

	if ok := core.IsJson(props); !ok {
		return core.NewError(core.InvalidProperties)
	}
	return nil
}

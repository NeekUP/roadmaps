package tests

import (
	"roadmaps/core/usecases"
	"roadmaps/domain"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"strings"
	"testing"
)

func TestAddBookIsbn13(t *testing.T) {
	isbn := "978-1-10-769989-2"
	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})
	result, err := usecase.Do(infrastructure.NewContext(nil), isbn, "{}", domain.Book)

	if err != nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn, err.Error())
	}

	if result == nil {
		t.Errorf("Book not saved as source using isbn %s without error", isbn)
	}

	if result.NormalizedIdentifier != strings.ReplaceAll(isbn, "-", "") {
		t.Errorf("Normalized isbn not valid. Expected: %s, Does: %s", strings.ReplaceAll(isbn, "-", ""), result.NormalizedIdentifier)
	}

	if result.Id == 0 {
		t.Errorf("Book saved, but id not assigned")
	}
}

func TestAddBookIsbn10(t *testing.T) {
	isbn10 := "3-598-21500-2"
	isbn13 := "978-3-598-21500-1"

	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})
	result, err := usecase.Do(infrastructure.NewContext(nil), isbn10, "{}", domain.Book)

	if err != nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn10, err.Error())
	}

	if result == nil {
		t.Errorf("Book not saved as source using isbn %s without error", isbn10)
	}

	if result.NormalizedIdentifier != strings.ReplaceAll(isbn13, "-", "") {
		t.Errorf("Normalized isbn not valid. Expected: %s, Does: %s", strings.ReplaceAll(isbn10, "-", ""), result.NormalizedIdentifier)
	}

	if result.Id == 0 {
		t.Errorf("Book saved, but id not assigned")
	}
}

func TestAddBookTwiceWithSameResult(t *testing.T) {
	isbn := "978-3-598-21501-8"
	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})
	result1, err := usecase.Do(infrastructure.NewContext(nil), isbn, "{}", domain.Book)

	result2, err := usecase.Do(infrastructure.NewContext(nil), isbn, "{}", domain.Book)

	if err != nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn, err.Error())
	}

	if result1.Id != result2.Id {
		t.Errorf("Second result returns with defferent id. 1:%d 2:%d", result1.Id, result2.Id)
	}

	if result1.NormalizedIdentifier != result2.NormalizedIdentifier {
		t.Errorf("Second result returns with defferent NormalizedIdentifier. 1:%s 2:%s", result1.NormalizedIdentifier, result2.NormalizedIdentifier)
	}

	if result1.Identifier != result2.Identifier {
		t.Errorf("Second result returns with defferent Identifier. 1:%s 2:%s", result1.Identifier, result2.Identifier)
	}

	if result1.Title != result2.Title {
		t.Errorf("Second result returns with defferent Title. 1:%s 2:%s", result1.Title, result2.Title)
	}
}

func TestAddBookBadIsbn13(t *testing.T) {
	isbn := "978-1-10-769989-0"
	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})
	result, err := usecase.Do(infrastructure.NewContext(nil), isbn, "{}", domain.Book)

	if err == nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn, err.Error())
	}

	if result != nil {
		t.Errorf("Book saved as source using isbn %s with error", isbn)
	}
}

func TestAddBookBadIsbn10(t *testing.T) {
	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})

	isbnList := []struct {
		x string
	}{
		{"3-598-21501-1"},
		{"3-598-xx-2"},
		{"3-598-2150x-1"},
		{""},
		{"dsdsd"},
		{"w-www-wwwww-w"},
		{"?-???-?????-?"},
		{" -   -     - "},
		{"          "},
		{"   "},
	}

	for _, isbn := range isbnList {
		result, err := usecase.Do(infrastructure.NewContext(nil), isbn.x, "{}", domain.Book)

		if err == nil {
			t.Errorf("Book not saved as source using isbn %s with error %s", isbn.x, err.Error())
		}

		if result != nil {
			t.Errorf("Book saved as source using isbn %s with error", isbn.x)
		}
	}
}

func TestAddLinkSuccess(t *testing.T) {
	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})

	linkList := []struct {
		url  string
		nUrl string
	}{
		{"http://ya.ru/", "ya.ru"},
		{"https://ya.ru", "ya.ru"},
		{"http://YA.Ru/", "ya.ru"},
		{"http://www.ya.ru/", "ya.ru"},
		{"https://stackoverflow.com/jobs?so_medium=StackOverflow&so_source=SiteNav", "stackoverflow.com/jobs?so_medium=stackoverflow&so_source=sitenav"},
		{"http://дом.рф/", "дом.рф"},
		{"http://xn--d1aqf.xn--p1ai/", "xn--d1aqf.xn--p1ai"},
	}

	for _, link := range linkList {
		result, err := usecase.Do(infrastructure.NewContext(nil), link.url, "{}", domain.Article)

		if err != nil {
			t.Errorf("Article not saved as source using url %s with error %s", link.url, err.Error())
		}

		if result == nil {
			t.Errorf("Article is null after saving source using link %s with error", link.url)
		}

		if result.NormalizedIdentifier != link.nUrl {
			t.Errorf("Article after saving using link %s contains not expected normalized url. Expected: %s Does: %s", link.url, link.nUrl, result.NormalizedIdentifier)
		}

		if result.Title == "" {
			t.Errorf("Title not defined for %s", link.url)
		}
	}
}

func TestTwitterSummary(t *testing.T) {
	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})

	linkList := []struct {
		url   string
		nUrl  string
		title string
	}{
		{"https://github.com/golang/go/issues/23669",
			"github.com/golang/go/issues/23669",
			"net/url: URL.String URL encodes a valid IDN domain · Issue #23669 · golang/go"},
	}

	for _, link := range linkList {
		result, err := usecase.Do(infrastructure.NewContext(nil), link.url, "{}", domain.Article)

		if err != nil {
			t.Errorf("Article not saved as source using url %s with error %s", link.url, err.Error())
		}

		if result == nil {
			t.Errorf("Article is null after saving source using link %s with error", link.url)
		}

		if result.NormalizedIdentifier != link.nUrl {
			t.Errorf("Article after saving using link %s contains not expected normalized url. Expected: %s Does: %s", link.url, link.nUrl, result.NormalizedIdentifier)
		}

		if result.Title != link.title {
			t.Errorf("Title not expected for %s: %s", link.url, result.Title)
		}

		if result.Img == "" {
			t.Errorf("Img not expected for %s: %s", link.url, result.Img)
		}
	}
}

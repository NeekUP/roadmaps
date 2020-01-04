package tests

import (
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"strings"
	"testing"
)

func TestAddBookIsbn13(t *testing.T) {
	u := registerUser("TestAddBookIsbn13", "TestAddBookIsbn13@w.ww", "TestAddBookIsbn13")
	if u != nil {
		defer DeleteUser(u.Id)
	} else {
		t.Error("User not registered")
		return
	}
	isbn := "978-1-10-769989-2"
	usecase := usecases.NewAddSource(db.NewSourceRepository(DB), log, &fakeImageManager{})
	source, err := usecase.Do(newContext(u), isbn, make(map[string]string), domain.Book)

	if err != nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn, err.Error())
	} else {
		defer DeleteSource(source.Id)
	}

	if source == nil {
		t.Errorf("Book not saved as source using isbn %s without error", isbn)
	}

	if source.NormalizedIdentifier != strings.ReplaceAll(isbn, "-", "") {
		t.Errorf("Normalized isbn not valid. Expected: %s, Does: %s", strings.ReplaceAll(isbn, "-", ""), source.NormalizedIdentifier)
	}

	if source.Id == 0 {
		t.Errorf("Book saved, but id not assigned")
	}
}

func TestAddBookIsbn10(t *testing.T) {
	u := registerUser("TestAddBookIsbn10", "TestAddBookIsbn10@w.ww", "TestAddBookIsbn10")
	if u != nil {
		defer DeleteUser(u.Id)
	}
	isbn10 := "3-598-21500-2"
	isbn13 := "978-3-598-21500-1"

	usecase := usecases.NewAddSource(db.NewSourceRepository(DB), log, &fakeImageManager{})
	source, err := usecase.Do(newContext(u), isbn10, make(map[string]string), domain.Book)

	if err != nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn10, err.Error())
	} else {
		defer DeleteSource(source.Id)
	}

	if source == nil {
		t.Errorf("Book not saved as source using isbn %s without error", isbn10)
	}

	if source.NormalizedIdentifier != strings.ReplaceAll(isbn13, "-", "") {
		t.Errorf("Normalized isbn not valid. Expected: %s, Does: %s", strings.ReplaceAll(isbn10, "-", ""), source.NormalizedIdentifier)
	}

	if source.Id == 0 {
		t.Errorf("Book saved, but id not assigned")
	}
}

func TestAddBookTwiceWithSameResult(t *testing.T) {
	u := registerUser("TestAddBookTwiceWithSameResult", "TestAddBookTwiceWithSameResult@w.ww", "TestAddBookTwiceWithSameResult")
	if u != nil {
		defer DeleteUser(u.Id)
	}
	isbn := "978-3-598-21501-8"
	usecase := usecases.NewAddSource(db.NewSourceRepository(DB), log, &fakeImageManager{})
	sourceOne, err := usecase.Do(newContext(u), isbn, make(map[string]string), domain.Book)
	if err != nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn, err.Error())
	} else {
		defer DeleteSource(sourceOne.Id)
	}

	sourceTwo, err := usecase.Do(infrastructure.NewContext(nil), isbn, make(map[string]string), domain.Book)
	if err != nil {
		t.Errorf("Book not saved as source using isbn %s with error %s", isbn, err.Error())
	} else {
		defer DeleteSource(sourceTwo.Id)
	}

	if sourceOne.Id != sourceTwo.Id {
		t.Errorf("Second result returns with defferent id. 1:%d 2:%d", sourceOne.Id, sourceTwo.Id)
	}

	if sourceOne.NormalizedIdentifier != sourceTwo.NormalizedIdentifier {
		t.Errorf("Second result returns with defferent NormalizedIdentifier. 1:%s 2:%s", sourceOne.NormalizedIdentifier, sourceTwo.NormalizedIdentifier)
	}

	if sourceOne.Identifier != sourceTwo.Identifier {
		t.Errorf("Second result returns with defferent Identifier. 1:%s 2:%s", sourceOne.Identifier, sourceTwo.Identifier)
	}

	if sourceOne.Title != sourceTwo.Title {
		t.Errorf("Second result returns with defferent Title. 1:%s 2:%s", sourceOne.Title, sourceTwo.Title)
	}
}

func TestAddBookBadIsbn13(t *testing.T) {
	u := registerUser("TestAddBookBadIsbn13", "TestAddBookBadIsbn13@w.ww", "TestAddBookBadIsbn13")
	if u != nil {
		defer DeleteUser(u.Id)
	}
	isbn := "978-1-10-769989-0"
	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})
	source, err := usecase.Do(newContext(u), isbn, make(map[string]string), domain.Book)

	if err == nil {
		defer DeleteSource(source.Id)
		t.Errorf("Book saved as source using isbn %s without error ", isbn)
	}

	if source != nil {
		t.Errorf("Book saved as source using isbn %s with error", isbn)
	}
}

func TestAddBookBadIsbn10(t *testing.T) {
	u := registerUser("TestAddBookBadIsbn10", "TestAddBookBadIsbn10@w.ww", "TestAddBookBadIsbn10")
	if u != nil {
		defer DeleteUser(u.Id)
	}
	usecase := usecases.NewAddSource(db.NewSourceRepository(DB), log, &fakeImageManager{})

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
		result, err := usecase.Do(newContext(u), isbn.x, make(map[string]string), domain.Book)

		if err == nil {
			defer DeleteSource(result.Id)
			t.Errorf("Book saved as source using isbn %s without error ", isbn.x)
		}

		if result != nil {
			t.Errorf("Book saved as source using isbn %s with error", isbn.x)
		}
	}
}

func TestAddLinkSuccess(t *testing.T) {
	u := registerUser("TestAddLinkSuccess", "TestAddLinkSuccess@w.ww", "TestAddLinkSuccess")
	if u != nil {
		defer DeleteUser(u.Id)
	}
	usecase := usecases.NewAddSource(db.NewSourceRepository(DB), log, &fakeImageManager{})

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
		result, err := usecase.Do(newContext(u), link.url, make(map[string]string), domain.Article)

		if err != nil {
			t.Errorf("Article not saved as source using url %s with error %s", link.url, err.Error())
		} else {
			defer DeleteSource(result.Id)
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

//func TestTwitterSummary(t *testing.T) {
//	usecase := usecases.NewAddSource(db.NewSourceRepository(nil), log, &fakeImageManager{})
//
//	linkList := []struct {
//		url   string
//		nUrl  string
//		title string
//	}{
//		{"https://github.com/golang/go/issues/23669",
//			"github.com/golang/go/issues/23669",
//			"net/url: URL.String URL encodes a valid IDN domain · Issue #23669 · golang/go"},
//	}
//
//	for _, link := range linkList {
//		result, err := usecase.Do(infrastructure.NewContext(nil), link.url, make(map[string]string), domain.Article)
//
//		if err != nil {
//			t.Errorf("Article not saved as source using url %s with error %s", link.url, err.Error())
//		}else{
//			defer DeleteSource(result.Id)
//		}
//
//		if result == nil {
//			t.Errorf("Article is null after saving source using link %s with error", link.url)
//		}
//
//		if result.NormalizedIdentifier != link.nUrl {
//			t.Errorf("Article after saving using link %s contains not expected normalized url. Expected: %s Does: %s", link.url, link.nUrl, result.NormalizedIdentifier)
//		}
//
//		if result.Title != link.title {
//			t.Errorf("Title not expected for %s: %s", link.url, result.Title)
//		}
//
//		if result.Img == "" {
//			t.Errorf("Img not expected for %s: %s", link.url, result.Img)
//		}
//	}
//}

package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"roadmaps/core"
	"roadmaps/domain"
	"strings"

	"golang.org/x/net/idna"

	"github.com/google/uuid"

	"github.com/nfnt/resize"

	_ "image/gif"
	jpeg "image/jpeg"
	_ "image/png"

	"github.com/PuerkitoBio/goquery"
	"github.com/moraes/isbn"
)

const (
	googleApiUrl      = "https://www.googleapis.com/books/v1/volumes?q=%s+isbn&fields=items/volumeInfo(title,subtitle,authors,description,industryIdentifiers,imageLinks)"
	openLibraryApiUrl = "http://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data"
)

type AddSource interface {
	Do(ctx core.ReqContext, identifier string, props map[string]string, sourceType domain.SourceType) (*domain.Source, error)
}

func NewAddSource(sr core.SourceRepository, log core.AppLogger, imgSaver core.ImageManager) AddSource {
	return &addSource{Repo: sr, Log: log, ImageManager: imgSaver}
}

type addSource struct {
	Repo         core.SourceRepository
	Log          core.AppLogger
	ImageManager core.ImageManager
}

type webPageSummary struct {
	Title string
	Img   image.Image
	Desc  string
}

type bookSummary struct {
	Title   string
	Isbn10  string
	Isbn13  string
	Img     image.Image
	Authors []string
	Desc    string
}

// Должен возвращать или уже созданный ранее или новый объект
func (this *addSource) Do(ctx core.ReqContext, identifier string, props map[string]string, sourceType domain.SourceType) (*domain.Source, error) {
	appErr := this.validate(identifier, props, sourceType)
	if props == nil {
		props = map[string]string{}
	}
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	s := &domain.Source{
		Identifier: identifier,
		Type:       sourceType}

	// Cast identifier to unified representation
	// link without protocol and www...
	// isbn to isbn13 format
	var err error
	s.NormalizedIdentifier, err = this.normalizeIdentifier(sourceType, identifier)
	if err != nil {
		this.Log.Errorw("Identifier contain not valid value",
			"ReqId", ctx.ReqId(),
			"Error", err,
			"Identifier", identifier)
		return nil, core.ValidationError(map[string]string{"identifier": core.InvalidFormat.String()})
	}

	// Find exists source by normalized identifier
	source := this.Repo.FindByIdentifier(s.NormalizedIdentifier)
	if source != nil {
		return source, nil
	}

	var img image.Image

	// Fetch source summary
	switch sourceType {
	case domain.Book:
		bookMeta, err := this.getBookMeta(s.NormalizedIdentifier)
		if err != nil {
			this.Log.Errorw("Book summary not parsed",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
				"Identifier", identifier)
			return nil, core.ValidationError(map[string]string{"identifier": core.SourceNotFound.String()})
		}

		s.Desc = bookMeta.Desc
		s.Title = bookMeta.Title
		img = bookMeta.Img
		props["isbn10"] = bookMeta.Isbn10
		props["isbn13"] = bookMeta.Isbn13
		props["authors"] = strings.Join(bookMeta.Authors, ", ")
		break
	case domain.Audio, domain.Video, domain.Article:
		pageMeta, err := this.getWebPageMeta(identifier)
		if err != nil {
			this.Log.Errorw("Fail to get page summary",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
				"Identifier", identifier)
			return nil, core.ValidationError(map[string]string{"identifier": core.SourceNotFound.String()})
		}

		s.Title = pageMeta.Title
		s.Desc = pageMeta.Desc
		img = pageMeta.Img
	default:
		return nil, core.ValidationError(map[string]string{"type": core.InvalidSourceType.String()})
	}

	// Resize and save image
	if img != nil {
		s.Img, err = this.resizeAndSaveImage(img)

		if err != nil {
			s.Img = ""
			this.Log.Errorw("Fail to save image",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
				"Identifier", identifier)
		}
	}

	p, err := json.Marshal(props)
	if err != nil {
		return nil, core.NewError(core.InvalidRequest)
	}

	s.Properties = string(p)
	s = this.Repo.GetOrAddByIdentifier(s)
	return s, nil
}

func (this *addSource) resizeAndSaveImage(img image.Image) (string, error) {
	resized, err := this.resizeImage(200, 0, img)
	if err != nil {
		return "", err
	}

	name := this.generateFileName("jpg")
	err = this.ImageManager.Save(resized, name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (this *addSource) normalizeIdentifier(sourceType domain.SourceType, identifier string) (string, error) {
	switch sourceType {
	case domain.Book:
		return this.getBookIdentifier(identifier)
	case domain.Audio, domain.Video, domain.Article:
		return this.getLinkIdentifier(identifier)
	default:
		return "", core.NewError(core.InvalidSourceType)
	}
}

func (this *addSource) getImage(uri string) (image.Image, error) {
	img, err := this.getImageByUrl(uri)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (this *addSource) getWebPageMeta(uri string) (*webPageSummary, error) {

	if this.isIDN(uri) {
		uri = this.decodeIDN(uri)
	}

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Http status code not OK: %d, %s, %s", res.StatusCode, res.Status, uri)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	summary := this.getTwitterMeta(doc)
	if summary != nil {
		return summary, nil
	}

	summary = this.getOpenGraphMeta(doc)
	if summary != nil {
		return summary, nil
	}

	summary = this.getRawHtmlMeta(doc)
	if summary != nil {
		return summary, nil
	}

	return nil, core.NewError(core.InaccessibleWebPage)
}

func (this *addSource) getBookMeta(isbn13 string) (*bookSummary, error) {
	summary, err := this.getBookMetaFromGoogle(isbn13)
	if err == nil {
		return summary, err
	}

	return this.getBookMetaFromOpenLibrary(isbn13)
}

func (this *addSource) getBookMetaFromGoogle(isbn13 string) (*bookSummary, error) {
	uri := fmt.Sprintf(googleApiUrl, isbn13)
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Google Book api return status [%s]", res.Status)
	}

	decoder := json.NewDecoder(res.Body)
	data := new(googleBooksSearch)
	err = decoder.Decode(data)

	if err != nil {
		return nil, err
	}

	// find suitable book info
	var summary *bookSummary
	for i := 0; i < len(data.Items); i++ {
		item := data.Items[i]
		if item.VolumeInfo != nil && item.VolumeInfo.IndustryIdentifiers != nil {
			for j := 0; j < len(item.VolumeInfo.IndustryIdentifiers); j++ {
				id := item.VolumeInfo.IndustryIdentifiers[j]
				if id.Identifier == isbn13 {
					summary = &bookSummary{
						Isbn13:  id.Identifier,
						Title:   item.VolumeInfo.Title,
						Authors: item.VolumeInfo.Authors,
						Desc:    item.VolumeInfo.Description,
					}

					if item.VolumeInfo.ImageLinks.Thumbnail != "" {
						summary.Img, err = this.getImage(item.VolumeInfo.ImageLinks.Thumbnail)
						if err == nil {
							this.Log.Errorw("Fail to download image from",
								"Url", item.VolumeInfo.ImageLinks.Thumbnail)
						}
					}
					break
				}
			}

			if summary != nil {
				for k := 0; k < len(item.VolumeInfo.IndustryIdentifiers); k++ {
					id := item.VolumeInfo.IndustryIdentifiers[k]
					if id.TypeName == "ISBN_10" {
						summary.Isbn10 = id.Identifier
					}
				}
				break
			}
		}
	}

	if summary != nil {
		return summary, nil
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return nil, fmt.Errorf("Fail to parse book info from response: %s", buf.String())
}

func (this *addSource) getBookMetaFromOpenLibrary(isbn13 string) (*bookSummary, error) {
	uri := fmt.Sprintf(openLibraryApiUrl, isbn13)
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("OpenLibrary Book api return status [%s]", res.Status)
	}

	decoder := json.NewDecoder(res.Body)
	data := new(map[string]openLibraryBook)
	err = decoder.Decode(data)

	if err != nil {
		return nil, err
	}

	var summary *bookSummary
	for _, value := range *data {
		authors := []string{}
		for _, author := range value.Authors {
			authors = append(authors, author.Name)
		}
		summary = &bookSummary{
			Isbn13:  value.Identifiers.Isbn13[0],
			Isbn10:  value.Identifiers.Isbn10[0],
			Title:   value.Title,
			Authors: authors,
			Desc:    "",
		}

		if value.Cover.Large != "" {
			summary.Img, err = this.getImage(value.Cover.Large)
			if err != nil {
				this.Log.Errorw("Fail to download image",
					"Url", value.Cover.Large,
					"Error", err.Error())
			}
		}
		break
	}

	if summary != nil {
		return summary, nil
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return nil, fmt.Errorf("Fail to parse book info from response: %s", buf.String())
}

func (this *addSource) decodeIDN(uri string) string {

	u, err := url.Parse(uri)

	unicodeHost, err := idna.Punycode.ToUnicode(u.Host)
	if err != nil {
		return uri
	}
	u.Host = unicodeHost
	uri, _ = url.PathUnescape(u.String())

	return uri
}

func (this *addSource) isIDN(uri string) bool {
	return strings.Contains(uri, "://xn--")
}

func (this *addSource) getTwitterMeta(doc *goquery.Document) *webPageSummary {

	summarycard := doc.Find("meta[name=twitter\\:card]")
	if summarycard.Length() == 0 {
		return nil
	}

	return this.getWebPageSummaryFromMeta(doc.Find("meta[name^=twitter\\:]"), "name", "content")
}

func (this *addSource) getOpenGraphMeta(doc *goquery.Document) *webPageSummary {

	selections := doc.Find("meta[property^=og\\:]")
	if selections.Length() == 0 {
		return nil
	}

	return this.getWebPageSummaryFromMeta(selections, "property", "content")
}

func (this *addSource) getRawHtmlMeta(doc *goquery.Document) *webPageSummary {
	meta := new(webPageSummary)
	titles := doc.Find("title")
	if titles.Length() != 1 {
		return nil
	}

	meta.Title = titles.First().Text()
	return meta
}

func (this *addSource) getWebPageSummaryFromMeta(selections *goquery.Selection, keyAttr, valueAttr string) *webPageSummary {
	meta := new(webPageSummary)
	selections.Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr(keyAttr); ok {

			contentType := val[strings.Index(val, ":")+1:]
			content := s.AttrOr(valueAttr, "")

			switch contentType {
			case "title":
				meta.Title = content
				break
			case "description":
				meta.Desc = content
				break
			case "image", "image:src":
				if meta.Img == nil {
					meta.Img, _ = this.getImage(content)
					if meta.Img == nil {
						this.Log.Errorw("Fail to download image from",
							"Url", content)
					}
				}
				break
			}
		}
	})

	if meta.Title == "" {
		return nil
	}

	return meta
}

func (this *addSource) generateFileName(extention string) string {
	return uuid.New().String() + "." + extention
}

func (this *addSource) getImageByUrl(uri string) (image.Image, error) {
	// check link without protocol
	u, err := url.Parse(uri)
	if err != nil {
		this.Log.Errorw("Fail to parse image url", "uri", uri, "error", err.Error())
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "https"
		uri = u.String()
	}

	response, err := http.Get(uri)
	if err != nil {
		this.Log.Errorw("Fail to download image", "uri", uri, "error", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	image, _, err := image.Decode(response.Body)
	if err != nil {
		this.Log.Errorw("Fail to decode image", "uri", uri, "error", err.Error())
		return nil, err
	}

	return image, nil
}

func (this *addSource) resizeImage(w, h uint, image image.Image) ([]byte, error) {

	newImage := resize.Resize(w, h, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, newImage, nil)
	if err != nil {
		this.Log.Errorw("Fail to resize image", "error", err.Error())
		return nil, err
	}

	return buf.Bytes(), nil
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

func (this *addSource) getBookIdentifier(identifier string) (string, error) {
	b := strings.Replace(identifier, "-", "", -1)
	if !isbn.Validate(b) {
		return "", fmt.Errorf("Isbn is not valid: %s", b)
	}

	if len(b) == 10 {
		return isbn.To13(b)
	} else {
		identifier = b
	}
	return identifier, nil
}

func (this *addSource) validate(identifier string, props map[string]string, sourceType domain.SourceType) *core.AppError {

	errors := make(map[string]string)
	if sourceType != domain.Book {
		u, err := url.Parse(identifier)

		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			errors["identifier"] = core.InvalidUrl.String()
		}

	} else {
		b := strings.Replace(identifier, "-", "", -1)
		if !isbn.Validate(b) {
			errors["identifier"] = core.InvalidISBN.String()
		}
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}

// google api types
type googleBooksSearch struct {
	TotalItems int    `json:"totalItems"`
	Items      []Item `json:"items"`
}

type industryIdentifier struct {
	TypeName   string `json:"type"`
	Identifier string `json:"identifier"`
}

type imageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

type VolumeInfo struct {
	Title               string               `json:"title"`
	Authors             []string             `json:"authors"`
	Oublisher           string               `json:"publisher"`
	Description         string               `json:"description"`
	IndustryIdentifiers []industryIdentifier `json:"industryIdentifiers"`
	ImageLinks          imageLinks           `json:"imageLinks"`
	Language            string               `json:"language"`
}

type Item struct {
	VolumeInfo *VolumeInfo `json:"volumeInfo"`
}

// Open library api types

type openLibraryBook struct {
	Title       string      `json:"title"`
	Cover       cover       `json:"cover"`
	Authors     []author    `json:"authors"`
	Identifiers identifiers `json:"identifiers"`
}

type cover struct {
	Large string `json:"large"`
}

type author struct {
	Name string `json:"name"`
}

type identifiers struct {
	Isbn13 []string `json:"isbn_13"`
	Isbn10 []string `json:"isbn_10"`
}

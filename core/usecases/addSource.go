package usecases

import (
	"bytes"
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

type AddSource interface {
	Do(ctx core.ReqContext, identifier, props string, sourceType domain.SourceType) (*domain.Source, error)
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
	Title   string
	Img     []byte
	ImgName string
	Desc    string
}

// Должен возвращать или уже созданный ранее или новый объект
func (this *addSource) Do(ctx core.ReqContext, identifier, props string, sourceType domain.SourceType) (*domain.Source, error) {
	err := this.validate(identifier, props, sourceType)
	if err != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return nil, err
	}

	s := &domain.Source{
		Identifier: identifier,
		Type:       sourceType,
		Properties: props}

	// Get normalized identifier
	switch sourceType {
	case domain.Book:
		s.NormalizedIdentifier = this.getBookIdentifier(identifier)
		break
	case domain.Audio, domain.Video, domain.Article:
		s.NormalizedIdentifier, _ = this.getLinkIdentifier(identifier)
		break
	default:
		return nil, core.NewError(core.InvalidSourceType)
	}

	if s.NormalizedIdentifier == "" {
		this.Log.Errorw("Normalized identifier is empty",
			"ReqId", ctx.ReqId(),
			"Identifier", identifier)
		return nil, core.NewError(core.InvalidRequest)
	}

	// Find exists data
	source := this.Repo.FindByIdentifier(s.NormalizedIdentifier)
	if source != nil {
		return source, nil
	}

	// Fetch source summary
	switch sourceType {
	case domain.Book:
		// TODO: get preview
		break
	case domain.Audio, domain.Video, domain.Article:
		pageMeta, err := this.getWebPageMeta(identifier)
		if err != nil {
			this.Log.Errorw("Fail to get page summary",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
				"Identifier", identifier)
			return nil, err
		}
		s.Title = pageMeta.Title
		s.Desc = pageMeta.Desc

		if len(pageMeta.Img) > 0 {
			err = this.ImageManager.Save(pageMeta.Img, pageMeta.ImgName)
			if err != nil {
				s.Img = ""
				this.Log.Errorw("Fail to save image",
					"ReqId", ctx.ReqId(),
					"Error", err.Error(),
					"Identifier", identifier)
			} else {
				s.Img = pageMeta.ImgName
			}
		}
	default:
		return nil, core.NewError(core.InvalidSourceType)
	}

	s = this.Repo.GetOrAddByIdentifier(s)
	return s, nil
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
				if len(meta.Img) == 0 {
					img, err := this.getImageByUrl(content)
					if err != nil {
						break
					}
					meta.Img, _ = this.resizeImage(160, 0, img)
					meta.ImgName = this.generateFileName("jpg")
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

func (this *addSource) getBookIdentifier(identifier string) string {
	b := strings.Replace(identifier, "-", "", -1)
	if len(b) == 10 {
		identifier, _ = isbn.To13(b)
	} else {
		identifier = b
	}
	return identifier
}

func (this *addSource) validate(identifier, props string, sourceType domain.SourceType) error {

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

	if ok := core.IsJson(props); !ok {
		return core.NewError(core.InvalidProperties)
	}
	return nil
}

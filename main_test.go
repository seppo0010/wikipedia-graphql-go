package main

import "os"
import "testing"
import "github.com/seppo0010/wikipedia-go"

func TestPageStrings(t *testing.T) {
	t.Parallel()
	query := "{page(title:\"Argentina\"){id,title,content,html_content,summary}}"
	result := executeQuery(query, schema)
	if result.Data.(map[string]interface{})["page"].(map[string]interface{})["id"].(string) != "1234" {
		t.Error("invalid page id")
		return
	}
	if result.Data.(map[string]interface{})["page"].(map[string]interface{})["title"].(string) != "Argentina" {
		t.Error("invalid page title")
		return
	}
	if result.Data.(map[string]interface{})["page"].(map[string]interface{})["content"].(string) != "Argentina is a country" {
		t.Error("invalid page content")
		return
	}
	if result.Data.(map[string]interface{})["page"].(map[string]interface{})["html_content"].(string) != "<b>Argentina</b> is a country" {
		t.Error("invalid page html content")
		return
	}
	if result.Data.(map[string]interface{})["page"].(map[string]interface{})["summary"].(string) != "The country" {
		t.Error("invalid page summary")
		return
	}
}

func TestPageImages(t *testing.T) {
	t.Parallel()
	query := "{page(title:\"Argentina\"){images{url,title,description_url}}}"
	result := executeQuery(query, schema)
	images := result.Data.(map[string]interface{})["page"].(map[string]interface{})["images"].([]interface{})
	if len(images) != 2 {
		t.Error("wrong image count")
		return
	}
	if images[0].(map[string]interface{})["url"].(string) != "url" {
		t.Error("wrong url")
		return
	}
	if images[1].(map[string]interface{})["url"].(string) != "url2" {
		t.Error("wrong url2")
		return
	}
	if images[0].(map[string]interface{})["title"].(string) != "title" {
		t.Error("wrong title")
		return
	}
	if images[1].(map[string]interface{})["title"].(string) != "title2" {
		t.Error("wrong title2")
		return
	}
	if images[0].(map[string]interface{})["description_url"].(string) != "description" {
		t.Error("wrong description")
		return
	}
	if images[1].(map[string]interface{})["description_url"].(string) != "description2" {
		t.Error("wrong description2")
		return
	}
}

func TestMain(m *testing.M) {
	wiki = NewWikipediaMock()
	wiki.(*WikipediaMock).AddPage(&PageMock{
		id:          "1234",
		title:       "Argentina",
		content:     "Argentina is a country",
		htmlContent: "<b>Argentina</b> is a country",
		summary:     "The country",
		images: []wikipedia.ImageRequest{
			{Image: wikipedia.Image{Url: "url", Title: "title", DescriptionUrl: "description"}},
			{Image: wikipedia.Image{Url: "url2", Title: "title2", DescriptionUrl: "description2"}},
		},
	})
	os.Exit(m.Run())
}
func NewWikipediaMock() *WikipediaMock {
	return &WikipediaMock{
		PagesById:    make(map[string]*PageMock),
		PagesByTitle: make(map[string]*PageMock),
	}
}

type WikipediaMock struct {
	PagesById    map[string]*PageMock
	PagesByTitle map[string]*PageMock
}
type PageMock struct {
	id             string
	idErr          error
	title          string
	titleErr       error
	content        string
	contentErr     error
	htmlContent    string
	htmlContentErr error
	summary        string
	summaryErr     error
	images         []wikipedia.ImageRequest
	imagesErr      error
}

func (w *WikipediaMock) AddPage(page *PageMock) {
	if page.id == "" && page.title == "" {
		panic("Page needs to have either an id or a title")
	}
	if page.id != "" {
		w.PagesById[page.id] = page
	}
	if page.title != "" {
		w.PagesByTitle[page.title] = page
	}
}

func (w *WikipediaMock) Page(title string) wikipedia.Page {
	return w.PagesByTitle[title]
}

func (w *WikipediaMock) PageFromId(id string) wikipedia.Page {
	return w.PagesById[id]
}
func (w *WikipediaMock) GetBaseUrl() string {
	return ""
}
func (w *WikipediaMock) SetBaseUrl(baseUrl string) {
}
func (w *WikipediaMock) SetImagesResults(imagesResults string) {
}
func (w *WikipediaMock) SetLinksResults(linksResults string)                       {}
func (w *WikipediaMock) SetCategoriesResults(categoriesResults string)             {}
func (w *WikipediaMock) PreLanguageUrl() string                                    { return "" }
func (w *WikipediaMock) PostLanguageUrl() string                                   { return "" }
func (w *WikipediaMock) Language() string                                          { return "" }
func (w *WikipediaMock) SearchResults() int                                        { return 0 }
func (w *WikipediaMock) GetLanguages() (languages []wikipedia.Language, err error) { return nil, nil }
func (w *WikipediaMock) Search(query string) (results []string, err error)         { return nil, nil }
func (w *WikipediaMock) Geosearch(latitude float64, longitude float64, radius int) (results []string, err error) {
	return nil, nil
}
func (w *WikipediaMock) RandomCount(count uint) (results []string, err error) { return nil, nil }
func (w *WikipediaMock) Random() (string, error)                              { return "", nil }
func (w *WikipediaMock) ImagesResults() string                                { return "" }
func (w *WikipediaMock) LinksResults() string                                 { return "" }
func (w *WikipediaMock) CategoriesResults() string                            { return "" }

func (p *PageMock) Id() (pageId string, err error)           { return p.id, p.idErr }
func (p *PageMock) Title() (pageTitle string, err error)     { return p.title, p.titleErr }
func (p *PageMock) Content() (content string, err error)     { return p.content, p.contentErr }
func (p *PageMock) HtmlContent() (content string, err error) { return p.htmlContent, p.htmlContentErr }
func (p *PageMock) Summary() (summary string, err error)     { return p.summary, p.summaryErr }
func (p *PageMock) Images() <-chan wikipedia.ImageRequest {
	ch := make(chan wikipedia.ImageRequest)
	go func() {
		for _, im := range p.images {
			ch <- im
		}
		close(ch)
	}()
	return ch
}
func (p *PageMock) Extlinks() <-chan wikipedia.ReferenceRequest {
	ch := make(chan wikipedia.ReferenceRequest)
	defer close(ch)
	return ch
}
func (p *PageMock) Links() <-chan wikipedia.LinkRequest {
	ch := make(chan wikipedia.LinkRequest)
	defer close(ch)
	return ch
}
func (p *PageMock) Categories() <-chan wikipedia.CategoryRequest {
	ch := make(chan wikipedia.CategoryRequest)
	defer close(ch)
	return ch
}
func (p *PageMock) Sections() (titles []string, err error)                         { return nil, nil }
func (p *PageMock) SectionContent(title string) (sectionContent string, err error) { return "", nil }

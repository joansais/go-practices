package wiki

import (
	"fmt"
	"regexp"
	"html/template"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var (
	linkPattern = regexp.MustCompile(`\[[a-zA-Z0-9]+\]`)
)

type Page struct {
	Title string
	Body  string
}

func listPages() ([]string, error) {
	return pageStore.List()
}

func loadPage(title string) (*Page, error) {
	return pageStore.Load(title)
}

func (page *Page) save() error {
	return pageStore.Save(page)
}

func (page *Page) remove() error {
	return pageStore.Remove(page.Title)
}

func (page *Page) Render() template.HTML {
	linkedMarkdown := linkPattern.ReplaceAllFunc([]byte(page.Body), linkToHtml)
	unsafeHtml := blackfriday.MarkdownCommon(linkedMarkdown)
	return template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafeHtml))
}

func linkToHtml(link []byte) []byte {
	title := link[1 : len(link)-1] // remove [...] brackets
	return []byte(fmt.Sprintf("<a href=\"/view/%s\">%s</a>", title, title))
}

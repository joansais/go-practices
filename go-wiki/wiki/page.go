package wiki

import (
	"fmt"
	"os"
	"io/ioutil"
	"regexp"
	"strings"
	"html/template"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

const (
	FILE_SUFFIX = ".wiki"
)

var (
	storageDir = "data/pages/"
	linkPattern = regexp.MustCompile(`\[[a-zA-Z0-9]+\]`)
)

func SetStorageDir(path string) {
	storageDir = path
}

type Page struct {
	Title string
	Body  string
}

func listPages() (result []string, err error) {
	files, err := ioutil.ReadDir(storageDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() && strings.HasSuffix(file.Name(), FILE_SUFFIX) {
			title := strings.TrimSuffix(file.Name(), FILE_SUFFIX)
			result = append(result, title)
		}
	}
	return
}

func loadPage(title string) (*Page, error) {
	content, err := ioutil.ReadFile(getPageFilename(title))
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(content)}, nil
}

func (page *Page) save() error {
	return ioutil.WriteFile(getPageFilename(page.Title), []byte(page.Body), 0600)
}

func (page *Page) remove() error {
	return os.Remove(getPageFilename(page.Title))
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

func getPageFilename(title string) string {
	return storageDir + title + FILE_SUFFIX
}

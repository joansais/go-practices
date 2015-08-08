package wiki

import (
	"html/template"
	"io/ioutil"
)

const (
	STORAGE_DIR = "data/pages/"
)

type Page struct {
	Title string
	Body  template.HTML
}

func (page *Page) save() error {
	return ioutil.WriteFile(getFilename(page.Title), []byte(page.Body), 0600)
}

func loadPage(title string) (*Page, error) {
	content, err := ioutil.ReadFile(getFilename(title))
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: template.HTML(content)}, nil
}

func getFilename(title string) string {
	return STORAGE_DIR + title + ".wiki"
}

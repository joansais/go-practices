package wiki

import (
    "testing"
    "html/template"
)

func TestPageLoadSave(t *testing.T) {
	pageSaved := Page{Title: "TestPage", Body: template.HTML("This is a sample page.")}
	pageSaved.save()
	
	pageLoaded, err := loadPage(pageSaved.Title)
	if err != nil {
		t.Error(err)
	}
	if pageLoaded.Title != pageSaved.Title {
		t.Error("Expected %q, found %q", pageSaved.Title, pageLoaded.Title)
	}
	if string(pageLoaded.Body) != string(pageSaved.Body) {
		t.Error("Expected %q, found %q", string(pageSaved.Body), string(pageLoaded.Body))
	}
}


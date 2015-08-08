package wiki

import (
	"sort"
	"testing"
)

var (
	samplePage = Page{Title: "TmpTestPage", Body: "This is a sample page for testing purposes."}
)

func TestLoadSave(t *testing.T) {
	samplePage.save()
	defer samplePage.remove()

	pageLoaded, err := loadPage(samplePage.Title)
	if err != nil {
		t.Error(err)
		return
	}
	if pageLoaded.Title != samplePage.Title {
		t.Errorf("Expected %q, found %q", samplePage.Title, pageLoaded.Title)
		return
	}
	if pageLoaded.Body != samplePage.Body {
		t.Errorf("Expected %q, found %q", samplePage.Body, pageLoaded.Body)
		return
	}
}

func TestRemove(t *testing.T) {
	samplePage.save()
	defer samplePage.remove()

	_, err := loadPage(samplePage.Title)
	if err != nil {
		t.Error(err)
		return
	}

	err = samplePage.remove()
	if err != nil {
		t.Error(err)
		return
	}

	_, err = loadPage(samplePage.Title)
	if err == nil {
		t.Error("Page not removed as expected")
		return
	}
}

func TestListPages(t *testing.T) {
	pages := []Page{
		Page{Title: "TmpTestPage1", Body: "This is a sample page for testing purposes."},
		Page{Title: "TmpTestPage3", Body: "This is a sample page for testing purposes."},
		Page{Title: "TmpTestPage2", Body: "This is a sample page for testing purposes."}}

	for _, page := range pages {
		page.save()
	}

	defer func() {
		for _, page := range pages {
			page.remove()
		}
	}()

	expected := []string{}
	for _, page := range pages {
		expected = append(expected, page.Title)
	}
	sort.Strings(expected)

	found, err := listPages()
	if err != nil {
		t.Error(err)
		return
	}

	if len(found) != len(expected) {
		t.Errorf("Expected %d pages, found %d", len(expected), len(found))
		return
	}

	for i := range found {
		if expected[i] != found[i] {
			t.Errorf("Expected %q, found %q", expected[i], found[i])
			return
		}
	}
}

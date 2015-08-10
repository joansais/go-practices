package wiki

import (
	"testing"
	"sort"
	"io/ioutil"
	"os"
)

func setupPageStore() *diskStore {
	storePath, err := ioutil.TempDir("", "wikitest")
	if err != nil {
		panic(err)
	}
	return &diskStore{path: storePath}
}

func cleanPageStore(store *diskStore) {
	ids, _ := store.ListAll()
	for _, id := range ids {
		store.Delete(id)
	}
	os.Remove(store.path)
}

func TestStoreCreateRead(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)
	
	page := &Page{Title: "Sample Page", Body: "This is a sample page for testing purposes."}
	id, err := store.Create(page)
	if err != nil {
		t.Error(err)
		return
	}
	if page.Id != id {
		t.Errorf("diskStore.Create: expected %q, found %q", page.Id, id)
		return
	}

	pageRead, err := store.Read(id)
	if err != nil {
		t.Error(err)
		return
	}

	if pageRead.Id != page.Id {
		t.Errorf("diskStore.Read(%q): expected %q, found %q", id, page.Id, pageRead.Id)
		return
	}
	if pageRead.Title != page.Title {
		t.Errorf("diskStore.Read(%q): expected %q, found %q", id, page.Title, pageRead.Title)
		return
	}
	if pageRead.Body != page.Body {
		t.Errorf("diskStore.Read(%q): expected %q, found %q", id, page.Body, pageRead.Body)
		return
	}
}

func TestStoreUpdate(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)
	
	page := &Page{Title: "Sample Page", Body: "This is a sample page for testing purposes."}
	id, err := store.Create(page)
	if err != nil {
		t.Error(err)
		return
	}
	
	page.Title = "Modified Page Title"
	page.Body = "This is a modified page body."

	err = store.Update(page)
	if err != nil {
		t.Error(err)
		return
	}

	pageRead, err := store.Read(id)
	if err != nil {
		t.Error(err)
		return
	}
	if pageRead.Title != page.Title {
		t.Errorf("diskStore.Read(%q): expected %q, found %q", id, page.Title, pageRead.Title)
		return
	}
	if pageRead.Body != page.Body {
		t.Errorf("diskStore.Read(%q): expected %q, found %q", id, page.Body, pageRead.Body)
		return
	}
}

func TestStoreDelete(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)

	page := &Page{Title: "Sample Page", Body: "This is a sample page for testing purposes."}
	id, err := store.Create(page)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = store.Read(id)
	if err != nil {
		t.Error(err)
		return
	}

	err = store.Delete(id)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = store.Read(id)
	if err == nil {
		t.Errorf("diskStore.Delete(%q): page was not deleted", id)
		return
	}
}

func TestStoreListAll(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)

	pages := []*Page{
		&Page{Title: "Sample Page 1", Body: "This is a sample page for testing purposes."},
		&Page{Title: "Sample Page 2", Body: "This is a sample page for testing purposes."},
		&Page{Title: "Sample Page 3", Body: "This is a sample page for testing purposes."}}

	expected := make([]string, len(pages))
	for k, page := range pages {
		id, err := store.Create(page)
		if err != nil {
			t.Error(err)
			return
		}
		expected[k] = string(id)
	}
	sort.Strings(expected)

	found, err := store.ListAll()
	if err != nil {
		t.Error(err)
		return
	}

	if len(found) != len(expected) {
		t.Errorf("diskStore.ListAll: expected %d pages, found %d", len(expected), len(found))
		return
	}

	for i := range found {
		if expected[i] != string(found[i]) {
			t.Errorf("diskStore.ListAll: expected %q, found %q", expected[i], found[i])
			return
		}
	}
}

func TestStoreUnexistentPageError(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)
	
	id := PageId("unexistent")
	_, err := store.Read(id)
	if err == nil {
		t.Error("diskStore.Read: an error was expected")
		return
	}

	unexistentPageErr, ok := err.(UnexistentPageError)
	if !ok {
		t.Error("diskStore.Read: UnexistentPageError was expected")
		return
	}
	
	if unexistentPageErr.Id != id {
		t.Error("UnexistentPageError.Id: expected %q, found %q", id, unexistentPageErr.Id)
		return
	}
}


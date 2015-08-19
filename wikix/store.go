package wikix

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/hex"
	"crypto/rand"
	"strings"
	"errors"
	"fmt"
	. "github.com/joansais/go-practices/exception"
)

type PageStore interface {
	Create(*Page) PageId
	Read(PageId) *Page
	Update(*Page)
	Delete(PageId)
	ListAll() []PageId
	FindByTitle(string) PageId
}

const (
	PAGE_ID_LEN = 6  // in bytes
	FILE_SUFFIX = ".wiki"
)

type diskStore struct {
	path string
}

func NewDiskStore(path string) PageStore {
	return &diskStore{path: path}
}

func (store *diskStore) Create(page *Page) PageId {
	id := newPageId()
	page.Id = id
	store.writePageToFile(page)
	return id
}

func (store *diskStore) Read(id PageId) *Page {
	return store.readPageFromFile(id)
}

func (store *diskStore) Update(page *Page) {
	_ = store.readPageFromFile(page.Id)  // check that page exists
	store.writePageToFile(page)
}

func (store *diskStore) Delete(id PageId) {
	filename := store.getPageFilename(id)
	ThrowIf(os.Remove(filename))
}

func (store *diskStore) ListAll() (result []PageId) {
	files, err := ioutil.ReadDir(store.path)
	ThrowIf(err)

	for _, file := range files {
		if file.Mode().IsRegular() && strings.HasSuffix(file.Name(), FILE_SUFFIX) {
			id := strings.TrimSuffix(file.Name(), FILE_SUFFIX)
			result = append(result, PageId(id))
		}
	}
	return
}

// TODO: implement more efficiently
func (store *diskStore) FindByTitle(title string) PageId {
	for _, id := range store.ListAll() {
		page := store.Read(id)
		if title == page.Title {
			return page.Id
		}
	}
	return ""
}

func (store *diskStore) readPageFromFile(id PageId) *Page {
	filename := store.getPageFilename(id)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			Throw(UnexistentPageError{id})
		} else {
			Throw(err)
		}
	}

	var page Page
	err = json.Unmarshal(content, &page)
	if err != nil {
		Throw(CorruptedFileError{filename, err})
	}
	
	if id != page.Id {
		Throw(CorruptedFileError{filename, errors.New("inconsistent page id")})
	}
	
	return &page
}

func (store *diskStore) writePageToFile(page *Page) {
	content, err := json.Marshal(page)
	ThrowIf(err)
	filename := store.getPageFilename(page.Id)
	err = ioutil.WriteFile(filename, content, 0600)
	ThrowIf(err)
}

func (store *diskStore) getPageFilename(id PageId) string {
	return store.path + "/" + string(id) + FILE_SUFFIX
}

func newPageId() PageId {
	bytes := make([]byte, PAGE_ID_LEN)
	_, err := rand.Read(bytes)
	ThrowIf(err)
	return PageId(hex.EncodeToString(bytes))
}

type UnexistentPageError struct {
    Id PageId
}

func (err UnexistentPageError) Error() string {
    return fmt.Sprintf("unexistent page %q", err.Id)
}

type CorruptedFileError struct {
    Filename string
    Cause error
}

func (err CorruptedFileError) Error() string {
    return fmt.Sprintf("corrupted file %q: %s", err.Filename, err.Cause.Error())
}


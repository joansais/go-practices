package wiki

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/hex"
	"crypto/rand"
	"strings"
	"errors"
	"fmt"
)

type PageStore interface {
	Create(*Page) (PageId, error)
	Read(PageId) (*Page, error)
	Update(*Page) error
	Delete(PageId) error
	ListAll() ([]PageId, error)
	FindByTitle(string) (PageId, error)
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

func (store *diskStore) Create(page *Page) (PageId, error) {
	id, err := newPageId()
	if err != nil {
		return "", err
	}
	
	page.Id = id
	err = store.writePageToFile(page)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (store *diskStore) Read(id PageId) (*Page, error) {
	page, err := store.readPageFromFile(id)
	if err != nil {
		return nil, err
	}
	
	return page, nil
}

func (store *diskStore) Update(page *Page) error {
	_, err := store.readPageFromFile(page.Id)  // check that page exists
	if err != nil {
		return err
	}
	
	return store.writePageToFile(page)
}

func (store *diskStore) Delete(id PageId) error {
	filename := store.getPageFilename(id)
	return os.Remove(filename)
}

func (store *diskStore) ListAll() (result []PageId, err error) {
	files, err := ioutil.ReadDir(store.path)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() && strings.HasSuffix(file.Name(), FILE_SUFFIX) {
			id := strings.TrimSuffix(file.Name(), FILE_SUFFIX)
			result = append(result, PageId(id))
		}
	}
	return
}

// TODO: implement more efficiently
func (store *diskStore) FindByTitle(title string) (PageId, error) {
	ids, err := store.ListAll()
	if err != nil {
		return "", err
	}

	for _, id := range ids {
		page, err := store.Read(id)
		if err != nil {
			return "", err
		}
		if title == page.Title {
			return page.Id, nil
		}
	}

	return "", nil
}

func (store *diskStore) readPageFromFile(id PageId) (*Page, error) {
	filename := store.getPageFilename(id)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, UnexistentPageError{id}
		} else {
			return nil, err
		}
	}

	var page Page
	err = json.Unmarshal(content, &page)
	if err != nil {
		return nil, CorruptedFileError{filename, err}
	}
	
	if id != page.Id {
		return nil, CorruptedFileError{filename, errors.New("inconsistent page id")}
	}
	
	return &page, nil
}

func (store *diskStore) writePageToFile(page *Page) error {
	content, err := json.Marshal(page)
	if err != nil {
		return err
	}

	filename := store.getPageFilename(page.Id)
	return ioutil.WriteFile(filename, content, 0600)
}

func (store *diskStore) getPageFilename(id PageId) string {
	return store.path + "/" + string(id) + FILE_SUFFIX
}

func newPageId() (PageId, error) {
	bytes := make([]byte, PAGE_ID_LEN)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return PageId(hex.EncodeToString(bytes)), nil
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


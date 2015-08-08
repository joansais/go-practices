package wiki

import (
	"os"
	"io/ioutil"
	"strings"
)

type PageStore interface {
	List() (result []string, err error)
	Load(title string) (*Page, error)
	Save(page *Page) error
	Remove(title string) error
}

var (
	pageStore PageStore
)

func SetPageStore(store PageStore) {
	pageStore = store
}

const (
	FILE_SUFFIX = ".wiki"
)

type diskStore struct {
	path string
}

func NewDiskStore(path string) PageStore {
	return &diskStore{path: path}
}

func (store *diskStore) List() (result []string, err error) {
	files, err := ioutil.ReadDir(store.path)
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

func (store *diskStore) Load(title string) (*Page, error) {
	content, err := ioutil.ReadFile(store.getPageFilename(title))
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(content)}, nil
}

func (store *diskStore) Save(page *Page) error {
	return ioutil.WriteFile(store.getPageFilename(page.Title), []byte(page.Body), 0600)
}

func (store *diskStore) Remove(title string) error {
	return os.Remove(store.getPageFilename(title))
}

func (store *diskStore) getPageFilename(title string) string {
	return store.path + "/" + title + FILE_SUFFIX
}


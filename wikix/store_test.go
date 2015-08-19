package wikix

import (
	"testing"
	"sort"
	. "gopkg.in/check.v1"
	. "github.com/joansais/go-practices/exception"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&StoreSuite{})

type StoreSuite struct {
	store PageStore
}

func (s *StoreSuite) SetUpSuite(c *C) {
	s.store = &diskStore{path: c.MkDir()}
}

func (s *StoreSuite) TearDownTest(c *C) {
	for _, id := range s.store.ListAll() {
		s.store.Delete(id)
	}
}

func (s *StoreSuite) TestCreateRead(c *C) {
	page := &Page{Title: "Sample Page", Body: "This is a sample page for testing purposes."}
	id := s.store.Create(page)

	c.Assert(id, Equals, page.Id)
	c.Assert(page, DeepEquals, s.store.Read(id))
}

func (s *StoreSuite) TestUpdate(c *C) {
	page := &Page{Title: "Sample Page", Body: "This is a sample page for testing purposes."}
	id := s.store.Create(page)
	
	page.Title = "Modified Page Title"
	page.Body = "This is a modified page body."

	s.store.Update(page)
	c.Assert(page, DeepEquals, s.store.Read(id))
}

func (s *StoreSuite) TestDelete(c *C) {
	page := &Page{Title: "Sample Page", Body: "This is a sample page for testing purposes."}
	id := s.store.Create(page)
	s.store.Read(id)
	s.store.Delete(id)
	c.Assert(Try(func() { s.store.Read(id) }), DeepEquals, UnexistentPageError{id})
}

func (s *StoreSuite) TestListAll(c *C) {
	pages := []*Page{
		&Page{Title: "Sample Page 1", Body: "This is a sample page for testing purposes."},
		&Page{Title: "Sample Page 2", Body: "This is a sample page for testing purposes."},
		&Page{Title: "Sample Page 3", Body: "This is a sample page for testing purposes."}}

	expected := make([]string, len(pages))
	for k, page := range pages {
		id := s.store.Create(page)
		expected[k] = string(id)
	}
	sort.Strings(expected)

	obtained := s.store.ListAll()
	
	c.Assert(len(obtained), Equals, len(expected))
	for k := range expected {
		c.Assert(string(obtained[k]), Equals, expected[k])
	}
}

func (s *StoreSuite) TestFindByTitle(c *C) {
	pages := []*Page{
		&Page{Title: "Sample Page 1", Body: "This is a sample page for testing purposes."},
		&Page{Title: "Sample Page 2", Body: "This is a sample page for testing purposes."},
		&Page{Title: "Sample Page 3", Body: "This is a sample page for testing purposes."}}

	for _, page := range pages {
		s.store.Create(page)
	}

	for _, page := range pages {
		c.Assert(s.store.FindByTitle(page.Title), Equals, page.Id)
	}

	c.Assert(s.store.FindByTitle("unexistent"), Equals, PageId(""))
}

func (s *StoreSuite) TestUnexistentPageError(c *C) {
	id := PageId("unexistent")
	c.Assert(Try(func() { s.store.Read(id) }), DeepEquals, UnexistentPageError{id})
}


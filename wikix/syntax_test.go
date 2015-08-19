package wikix

import (
	"bytes"
	"encoding/xml"
	"fmt"
	. "gopkg.in/check.v1"
)

var _ = Suite(&SyntaxSuite{})

type SyntaxSuite struct {
	StoreSuite
}

func (s *SyntaxSuite) SetUpSuite(c *C) {
	s.StoreSuite.SetUpSuite(c)
}

func (s *SyntaxSuite) TearDownTest(c *C) {
	s.StoreSuite.TearDownTest(c)
}

func (s *SyntaxSuite) TestRegularMarkdown(c *C) {
	syntax := &markdownSyntax{s.store}

	body := "Some text with *markdown*, an [inline link](http://example1.net/), a [reference link] [1], and another [REF LINK][].\n" +
		"[1]: http://example2.net/\n" +
		"[REF LINK]: http://example3.net/ \"Optional title\" \n"

	obtained := syntax.BodyToEdit(body)
	expected := body
	c.Assert(obtained, Equals, expected)  // BodyToEdit

	obtained = syntax.EditToBody(obtained)
	expected = body
	c.Assert(obtained, Equals, expected)  // EditToBody

	obtained = string(syntax.BodyToHtml(body))
	expected = "<p>Some text with <em>markdown</em>, an <a href=\"http://example1.net/\" rel=\"nofollow\">inline link</a>, " +
		"a <a href=\"http://example2.net/\" rel=\"nofollow\">reference link</a>, " +
		"and another <a href=\"http://example3.net/\" title=\"Optional title\" rel=\"nofollow\">REF LINK</a>.</p>\n"
	c.Assert(obtained, Equals, expected)  // BodyToHtml
}

func (s *SyntaxSuite) TestPageLinkEditing(c *C) {
	syntax := &markdownSyntax{s.store}

	page1 := &Page{Title: "Page #1", Body: "Some text with *markdown*."}
	s.store.Create(page1)

	page2 := &Page{Title: "Page #2", Body: fmt.Sprintf("Some text referencing [%s][].", page1.Id)}
	s.store.Create(page2)

	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)}
	s.store.Create(page3)

	obtained := syntax.BodyToEdit(page1.Body)
	expected := page1.Body
	c.Assert(obtained, Equals, expected)  // BodyToEdit

	obtained = syntax.BodyToEdit(page2.Body)
	expected = fmt.Sprintf("Some text referencing [%s][].", page1.Title)
	c.Assert(obtained, Equals, expected)  // BodyToEdit

	obtained = syntax.EditToBody(obtained)
	expected = fmt.Sprintf("Some text referencing [%s][].", page1.Id)
	c.Assert(obtained, Equals, expected)  // EditToBody

	obtained = syntax.BodyToEdit(page3.Body)
	expected = fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Title, page1.Title)
	c.Assert(obtained, Equals, expected)  // BodyToEdit

	obtained = syntax.EditToBody(obtained)
	expected = fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)
	c.Assert(obtained, Equals, expected)  // EditToBody
}

func (s *SyntaxSuite) TestPageLinkRendering(c *C) {
	syntax := &markdownSyntax{s.store}

	page1 := &Page{Title: "Page #1", Body: "Some text with *markdown*."}
	s.store.Create(page1)

	page2 := &Page{Title: "Page #2", Body: fmt.Sprintf("Some text referencing [%s][].", page1.Id)}
	s.store.Create(page2)

	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)}
	s.store.Create(page3)

	obtained := string(syntax.BodyToHtml(page1.Body))
	expected := "<p>Some text with <em>markdown</em>.</p>\n"
	c.Assert(obtained, Equals, expected)  // BodyToHtml

	obtained = string(syntax.BodyToHtml(page2.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page1.Id, page1.Title, page1.Title)
	c.Assert(obtained, Equals, expected)  // BodyToHtml

	obtained = string(syntax.BodyToHtml(page3.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">the second page</a> " +
		"and <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page2.Id, page2.Title, page1.Id, page1.Title, page1.Title)
	c.Assert(obtained, Equals, expected)  // BodyToHtml
}

func (s *SyntaxSuite) TestPageLinkSpecialChars(c *C) {
	syntax := &markdownSyntax{s.store}

	page1 := &Page{Title: "Page #1 \"with quotes\"", Body: "Some text with *markdown*."}
	s.store.Create(page1)

	page2 := &Page{Title: "Page #2 <with brackets>", Body: fmt.Sprintf("Some text referencing [%s][].", page1.Id)}
	s.store.Create(page2)

	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)}
	s.store.Create(page3)

	obtained := string(syntax.BodyToHtml(page1.Body))
	expected := "<p>Some text with <em>markdown</em>.</p>\n"
	c.Assert(obtained, Equals, expected)  // BodyToHtml

	obtained = string(syntax.BodyToHtml(page2.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page1.Id, escapeXml(page1.Title), escapeXml(page1.Title))
	c.Assert(obtained, Equals, expected)  // BodyToHtml

	obtained = string(syntax.BodyToHtml(page3.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">the second page</a> " +
		 "and <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page2.Id, escapeXml(page2.Title), page1.Id, escapeXml(page1.Title), escapeXml(page1.Title))
	c.Assert(obtained, Equals, expected)  // BodyToHtml
}

func escapeXml(in string) string {
	var out bytes.Buffer
	err := xml.EscapeText(&out, []byte(in))
	if err != nil {
		panic("Could not escape XML: " + err.Error())
	}
	return out.String()
}

package wiki

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"
)

func TestSyntaxRegularMarkdown(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)
	syntax := &markdownSyntax{store}

	body := "Some text with *markdown*, an [inline link](http://example1.net/), a [reference link] [1], and another [REF LINK][].\n" +
		"[1]: http://example2.net/\n" +
		"[REF LINK]: http://example3.net/ \"Optional title\" \n"

	obtained := syntax.BodyToEdit(body)
	expected := body
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToEdit: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.EditToBody(obtained)
	expected = body
	if obtained != expected {
		t.Errorf("markdownSyntax.EditToBody: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = string(syntax.BodyToHtml(body))
	expected = "<p>Some text with <em>markdown</em>, an <a href=\"http://example1.net/\" rel=\"nofollow\">inline link</a>, " +
		"a <a href=\"http://example2.net/\" rel=\"nofollow\">reference link</a>, " +
		"and another <a href=\"http://example3.net/\" title=\"Optional title\" rel=\"nofollow\">REF LINK</a>.</p>\n"
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}
}

func TestSyntaxPageLinkEditing(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)
	syntax := &markdownSyntax{store}

	page1 := &Page{Title: "Page #1", Body: "Some text with *markdown*."}
	_, err := store.Create(page1)
	if err != nil {
		t.Error(err)
		return
	}

	page2 := &Page{Title: "Page #2", Body: fmt.Sprintf("Some text referencing [%s][].", page1.Id)}
	_, err = store.Create(page2)
	if err != nil {
		t.Error(err)
		return
	}

	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)}
	_, err = store.Create(page3)
	if err != nil {
		t.Error(err)
		return
	}

	obtained := syntax.BodyToEdit(page1.Body)
	expected := page1.Body
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToEdit: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.BodyToEdit(page2.Body)
	expected = fmt.Sprintf("Some text referencing [%s][].", page1.Title)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToEdit: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.EditToBody(obtained)
	expected = fmt.Sprintf("Some text referencing [%s][].", page1.Id)
	if obtained != expected {
		t.Errorf("markdownSyntax.EditToBody: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.BodyToEdit(page3.Body)
	expected = fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Title, page1.Title)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToEdit: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.EditToBody(obtained)
	expected = fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)
	if obtained != expected {
		t.Errorf("markdownSyntax.EditToBody: expected %q, obtained %q", expected, obtained)
		return
	}
}

func TestSyntaxPageLinkRendering(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)
	syntax := &markdownSyntax{store}

	page1 := &Page{Title: "Page #1", Body: "Some text with *markdown*."}
	_, err := store.Create(page1)
	if err != nil {
		t.Error(err)
		return
	}

	page2 := &Page{Title: "Page #2", Body: fmt.Sprintf("Some text referencing [%s][].", page1.Id)}
	_, err = store.Create(page2)
	if err != nil {
		t.Error(err)
		return
	}

	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)}
	_, err = store.Create(page3)
	if err != nil {
		t.Error(err)
		return
	}

	obtained := string(syntax.BodyToHtml(page1.Body))
	expected := "<p>Some text with <em>markdown</em>.</p>\n"
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = string(syntax.BodyToHtml(page2.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page1.Id, page1.Title, page1.Title)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = string(syntax.BodyToHtml(page3.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">the second page</a> " +
		"and <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page2.Id, page2.Title, page1.Id, page1.Title, page1.Title)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}
}

func TestSyntaxPageLinkSpecialChars(t *testing.T) {
	store := setupPageStore()
	defer cleanPageStore(store)
	syntax := &markdownSyntax{store}

	page1 := &Page{Title: "Page #1 \"with quotes\"", Body: "Some text with *markdown*."}
	_, err := store.Create(page1)
	if err != nil {
		t.Error(err)
		return
	}

	page2 := &Page{Title: "Page #2 <with brackets>", Body: fmt.Sprintf("Some text referencing [%s][].", page1.Id)}
	_, err = store.Create(page2)
	if err != nil {
		t.Error(err)
		return
	}

	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing [the second page][%s] and [%s] [].", page2.Id, page1.Id)}
	_, err = store.Create(page3)
	if err != nil {
		t.Error(err)
		return
	}

	obtained := string(syntax.BodyToHtml(page1.Body))
	expected := "<p>Some text with <em>markdown</em>.</p>\n"
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = string(syntax.BodyToHtml(page2.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page1.Id, escapeXml(page1.Title), escapeXml(page1.Title))
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = string(syntax.BodyToHtml(page3.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">the second page</a> " +
		 "and <a href=\"/view/%s\" title=\"%s\" rel=\"nofollow\">%s</a>.</p>\n", page2.Id, escapeXml(page2.Title), page1.Id, escapeXml(page1.Title), escapeXml(page1.Title))
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}
}

func escapeXml(in string) string {
	var out bytes.Buffer
	err := xml.EscapeText(&out, []byte(in))
	if err != nil {
		return in
	}
	return out.String()
}

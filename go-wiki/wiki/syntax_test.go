package wiki

import (
	"testing"
	"fmt"
)

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
	
	page2 := &Page{Title: "Page #2", Body: fmt.Sprintf("Some text referencing {%s}.", page1.Id)}
	_, err = store.Create(page2)
	if err != nil {
		t.Error(err)
		return
	}
	
	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing {the second page|%s}.", page2.Id)}
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
	expected = fmt.Sprintf("Some text referencing {%s}.", page1.Title)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToEdit: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.EditToBody(obtained)
	expected = fmt.Sprintf("Some text referencing {%s}.", page1.Id)
	if obtained != expected {
		t.Errorf("markdownSyntax.EditToBody: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.BodyToEdit(page3.Body)
	expected = fmt.Sprintf("Some text referencing {the second page|%s}.", page2.Title)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToEdit: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = syntax.EditToBody(obtained)
	expected = fmt.Sprintf("Some text referencing {the second page|%s}.", page2.Id)
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
	
	page2 := &Page{Title: "Page #2", Body: fmt.Sprintf("Some text referencing {%s}.", page1.Id)}
	_, err = store.Create(page2)
	if err != nil {
		t.Error(err)
		return
	}
	
	page3 := &Page{Title: "Page #3", Body: fmt.Sprintf("Some text referencing {the second page|%s}.", page2.Id)}
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
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" rel=\"nofollow\">%s</a>.</p>\n", page1.Id, page1.Title)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}

	obtained = string(syntax.BodyToHtml(page3.Body))
	expected = fmt.Sprintf("<p>Some text referencing <a href=\"/view/%s\" rel=\"nofollow\">the second page</a>.</p>\n", page2.Id)
	if obtained != expected {
		t.Errorf("markdownSyntax.BodyToHtml: expected %q, obtained %q", expected, obtained)
		return
	}
}

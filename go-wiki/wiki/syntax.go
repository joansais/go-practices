package wiki

import (
	"fmt"
	"regexp"
	"html/template"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"strings"
)

var (
	pageLinkPattern = regexp.MustCompile(`{.+}`)
)

type SyntaxHandler interface {
	BodyToEdit(body string) string
	EditToBody(edit string) string
	BodyToHtml(body string) template.HTML
}

type markdownSyntax struct {
	pageStore PageStore
}

func NewMarkdownSyntax(store PageStore) SyntaxHandler {
	return &markdownSyntax{store}
}

func (syntax *markdownSyntax) BodyToEdit(body string) string {
	return pageLinkPattern.ReplaceAllStringFunc(body, syntax.pageIdToTitle)
}

func (syntax *markdownSyntax) pageIdToTitle(str string) string {
	link := parseLink(str)

	page, err := syntax.pageStore.Read(PageId(link.ref))
	if err == nil {   // link references an existent page
		link.ref = page.Title
	}
	
	return link.String()
}

func (syntax *markdownSyntax) EditToBody(edit string) string {
	return pageLinkPattern.ReplaceAllStringFunc(edit, syntax.titleToPageId)
}

func (syntax *markdownSyntax) titleToPageId(str string) string {
	link := parseLink(str)

	pageId, err := syntax.pageStore.FindByTitle(link.ref)
	if err == nil {  // link references an existent page
		link.ref = string(pageId)
	}
	
	return link.String()
}

func (syntax *markdownSyntax) BodyToHtml(body string) template.HTML {
	unsafeHtml := string(blackfriday.MarkdownCommon([]byte(body)))
	unsafeHtml = pageLinkPattern.ReplaceAllStringFunc(unsafeHtml, syntax.pageIdToLink)
	return template.HTML(bluemonday.UGCPolicy().Sanitize(unsafeHtml))
}

func (syntax *markdownSyntax) pageIdToLink(str string) string {
	link := parseLink(str)

	var url, txt, tit string
	page, err := syntax.pageStore.Read(PageId(link.ref))
	if err == nil {   // link references an existent page
		url = "/view/" + link.ref
		if link.txt != "" {
			txt = link.txt
		} else {
			txt = page.Title
		}
		tit = "Go to page: " + page.Title
	} else {   // allow creating a page with that title
		url = "/create/?title=" + link.ref
		if link.txt != "" {
			txt = link.txt
		} else {
			txt = link.ref
		}
		tit = "Add page: " + link.ref
	}
	
	return fmt.Sprintf("<a href=%q title=%q>%s</a>", url, tit, txt)
}

type pageLink struct {
	txt string
	ref string
}

func parseLink(str string) *pageLink {
	inner := str[1 : len(str)-1] // remove {...} brackets
	items := strings.SplitN(inner, "|", 2)

	var txt, ref string
	if len(items) == 2 {
		txt = items[0]
		ref = items[1]
	} else {
		txt = ""
		ref = items[0]
	}

	return &pageLink{txt: strings.TrimSpace(txt), ref: strings.TrimSpace(ref)}
}

func (link pageLink) String() string {
	if link.txt != "" {
		return fmt.Sprintf("{%s|%s}", link.txt, link.ref)
	} else {
		return fmt.Sprintf("{%s}", link.ref)
	}
}

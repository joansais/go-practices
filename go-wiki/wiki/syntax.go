package wiki

import (
	"fmt"
	"regexp"
	"html/template"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

var (
	// FIXME: this regexp matches "[x] etc [y][z]", it should not!
	pageLinkPattern = regexp.MustCompile(`\[(.+)\]\[(.*)\]`)
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

func (syntax *markdownSyntax) pageIdToTitle(linkStr string) string {
	link := parseLink(linkStr)

	page, err := syntax.pageStore.Read(PageId(link.ref))
	if err == nil {   // link references an existent page
		link.ref = page.Title
	}
	
	return link.String()
}

func (syntax *markdownSyntax) EditToBody(edit string) string {
	return pageLinkPattern.ReplaceAllStringFunc(edit, syntax.titleToPageId)
}

func (syntax *markdownSyntax) titleToPageId(linkStr string) string {
	link := parseLink(linkStr)

	pageId, err := syntax.pageStore.FindByTitle(link.ref)
	if err == nil && pageId != "" {  // link references an existent page
		link.ref = string(pageId)
	}
	
	return link.String()
}

func (syntax *markdownSyntax) BodyToHtml(body string) template.HTML {
	renderer, options := syntax.markdownParams()
	unsafeHtml := string(blackfriday.MarkdownOptions([]byte(body), renderer, options))
	return template.HTML(sanitizePolicy().Sanitize(unsafeHtml))
}

// FIXME: commonHtmlFlags and commonExtensions should be exported by blackfriday
func (syntax *markdownSyntax) markdownParams() (blackfriday.Renderer, blackfriday.Options) {
	renderer := blackfriday.HtmlRenderer(blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES, "", "")
	
	options := blackfriday.Options{Extensions: blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS,
		ReferenceOverride: syntax.pageIdToLink }
	
	return renderer, options
}

func sanitizePolicy() *bluemonday.Policy {
	return bluemonday.UGCPolicy().AllowAttrs("title").OnElements("a");
}

func (syntax *markdownSyntax) pageIdToLink(reference string) (ref *blackfriday.Reference, overridden bool) {
	page, err := syntax.pageStore.Read(PageId(reference))
	if err == nil {  // reference to an existing page
		link := fmt.Sprintf("/view/%s", page.Id)
		ref = &blackfriday.Reference{Link: link, Title: page.Title, Text: page.Title}
		overridden = true
	} else {  // not referencing a page, or I/O error occurred
		ref = nil
		overridden = false
	}
	return
}

type pageLink struct {
	txt string
	ref string
}

func parseLink(linkStr string) *pageLink {
	submatches := pageLinkPattern.FindStringSubmatch(linkStr)
/*
fmt.Printf("linkStr: %s\n", linkStr)
for k, submatch := range submatches {
	fmt.Printf("submatches[%d]: %s\n", k, submatch)
}
*/
	var txt, ref string
	if submatches[2] != "" {
		txt = submatches[1]
		ref = submatches[2]
	} else {
		txt = ""
		ref = submatches[1]
	}

	return &pageLink{txt: strings.TrimSpace(txt), ref: strings.TrimSpace(ref)}
}

func (link pageLink) String() string {
	if link.txt != "" {
		return fmt.Sprintf("[%s][%s]", link.txt, link.ref)
	} else {
		return fmt.Sprintf("[%s][]", link.ref)
	}
}

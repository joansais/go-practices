package wiki

import (
	"net/http"
	"html/template"
	"regexp"
	"errors"
	"sort"
)

const (
	LIST_ENTRYPOINT_PATH = "/"
	VIEW_ENTRYPOINT_PATH = "/view/"
	CREATE_ENTRYPOINT_PATH = "/create/"
	EDIT_ENTRYPOINT_PATH = "/edit/"
	SAVE_ENTRYPOINT_PATH = "/save/"
	DELETE_ENTRYPOINT_PATH = "/delete/"
	HTML_TEMPLATE_FILES  = "/html/*.tmpl"
)

var (
	pageRequestPattern = regexp.MustCompile(`^/(view|edit|delete)/([a-zA-Z0-9]+)$`)
	pageStore PageStore
	syntaxHandler SyntaxHandler
	htmlTemplates *template.Template
)

func SetPageStore(store PageStore) {
	pageStore = store
}

func SetSyntaxHandler(syntax SyntaxHandler) {
	syntaxHandler = syntax
}

func SetAssetsDir(path string) {
	htmlTemplates = template.Must(template.ParseGlob(path + HTML_TEMPLATE_FILES))
}

type PageModel struct {
	Id			PageId
	Title		string
	BodyToEdit	string
	BodyAsHtml	template.HTML
}

type PageListModel []*PageModel

func (list PageListModel) Len() int           { return len(list) }
func (list PageListModel) Less(i, j int) bool { return list[i].Title < list[j].Title }
func (list PageListModel) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }

func handleList(res http.ResponseWriter, req *http.Request) {
	ids, err := pageStore.ListAll()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	pageList := make(PageListModel, len(ids))
	for k, id := range ids {
		page, err := pageStore.Read(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		pageList[k] = &PageModel{Id: id, Title: page.Title}
	}
	sort.Sort(pageList)

	err = htmlTemplates.ExecuteTemplate(res, "list", pageList)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleView(res http.ResponseWriter, req *http.Request) {
	id, err := getRequestedPageId(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := pageStore.Read(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	bodyAsHtml := syntaxHandler.BodyToHtml(page.Body)
	pageModel := &PageModel{Id: id, Title: page.Title, BodyAsHtml: bodyAsHtml}

	err = htmlTemplates.ExecuteTemplate(res, "view", pageModel)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCreate(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Query().Get("title")
	if title == "" {
		title = "New Page"
	}

	pageModel := &PageModel{Title: title}

	err := htmlTemplates.ExecuteTemplate(res, "create", pageModel)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEdit(res http.ResponseWriter, req *http.Request) {
	id, err := getRequestedPageId(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := pageStore.Read(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	bodyToEdit := syntaxHandler.BodyToEdit(page.Body)
	pageModel := &PageModel{Id: id, Title: page.Title, BodyToEdit: bodyToEdit}

	err = htmlTemplates.ExecuteTemplate(res, "edit", pageModel)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSave(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	id := PageId(req.Form.Get("id"))  // when creating a page, id == ""
	title := req.Form.Get("title")
	bodyFromEdit := req.Form.Get("body")
	body := syntaxHandler.EditToBody(bodyFromEdit)
	page := &Page{Id: id, Title: title, Body: body}

	if id == "" {
		id, err = pageStore.Create(page)
	} else {
		err = pageStore.Update(page)
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, VIEW_ENTRYPOINT_PATH+string(id), http.StatusFound)
}

func handleDelete(res http.ResponseWriter, req *http.Request) {
	id, err := getRequestedPageId(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = pageStore.Delete(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, LIST_ENTRYPOINT_PATH, http.StatusFound)
}

func getRequestedPageId(req *http.Request) (PageId, error) {
	submatches := pageRequestPattern.FindStringSubmatch(req.URL.Path)
	if submatches == nil {
		return "", errors.New("invalid page id")
	}
	return PageId(submatches[2]), nil
}

func RegisterServices() {
	http.HandleFunc(LIST_ENTRYPOINT_PATH, handleList)
	http.HandleFunc(VIEW_ENTRYPOINT_PATH, handleView)
	http.HandleFunc(CREATE_ENTRYPOINT_PATH, handleCreate)
	http.HandleFunc(EDIT_ENTRYPOINT_PATH, handleEdit)
	http.HandleFunc(SAVE_ENTRYPOINT_PATH, handleSave)
	http.HandleFunc(DELETE_ENTRYPOINT_PATH, handleDelete)
}

package wiki

import (
	"net/http"
	"html/template"
	"regexp"
	"errors"
	"sort"
	"log"
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
)

type Server struct {
	pageStore PageStore
	syntaxHandler SyntaxHandler
	htmlTemplates *template.Template
}

func NewServer(store PageStore, syntax SyntaxHandler, assetsDir string) *Server {
	return &Server{
		pageStore: store,
		syntaxHandler: syntax,
		htmlTemplates: template.Must(template.ParseGlob(assetsDir + HTML_TEMPLATE_FILES))}
}

func (server *Server) Start(addr string) error {
	http.HandleFunc(LIST_ENTRYPOINT_PATH, server.handleList)
	http.HandleFunc(VIEW_ENTRYPOINT_PATH, server.handleView)
	http.HandleFunc(CREATE_ENTRYPOINT_PATH, server.handleCreate)
	http.HandleFunc(EDIT_ENTRYPOINT_PATH, server.handleEdit)
	http.HandleFunc(SAVE_ENTRYPOINT_PATH, server.handleSave)
	http.HandleFunc(DELETE_ENTRYPOINT_PATH, server.handleDelete)
	return http.ListenAndServe(addr, nil)
}

func (server *Server) handleList(res http.ResponseWriter, req *http.Request) {
	ids, err := server.pageStore.ListAll()
	if err != nil {
		internalError(res, err)
		return
	}

	pageList := make(PageListModel, len(ids))
	for k, id := range ids {
		page, err := server.pageStore.Read(id)
		if err != nil {
			internalError(res, err)
			return
		}
		pageList[k] = &PageModel{Id: id, Title: page.Title}
	}
	sort.Sort(pageList)

	err = server.htmlTemplates.ExecuteTemplate(res, "list", pageList)
	if err != nil {
		internalError(res, err)
		return
	}
}

func (server *Server) handleView(res http.ResponseWriter, req *http.Request) {
	id, err := getRequestedPageId(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := server.pageStore.Read(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	bodyAsHtml := server.syntaxHandler.BodyToHtml(page.Body)
	pageModel := &PageModel{Id: id, Title: page.Title, BodyAsHtml: bodyAsHtml}

	err = server.htmlTemplates.ExecuteTemplate(res, "view", pageModel)
	if err != nil {
		internalError(res, err)
		return
	}
}

func (server *Server) handleCreate(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Query().Get("title")
	if title == "" {
		title = "New Page"
	}

	pageModel := &PageModel{Title: title}

	err := server.htmlTemplates.ExecuteTemplate(res, "create", pageModel)
	if err != nil {
		internalError(res, err)
		return
	}
}

func (server *Server) handleEdit(res http.ResponseWriter, req *http.Request) {
	id, err := getRequestedPageId(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := server.pageStore.Read(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	bodyToEdit := server.syntaxHandler.BodyToEdit(page.Body)
	pageModel := &PageModel{Id: id, Title: page.Title, BodyToEdit: bodyToEdit}

	err = server.htmlTemplates.ExecuteTemplate(res, "edit", pageModel)
	if err != nil {
		internalError(res, err)
		return
	}
}

func (server *Server) handleSave(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	id := PageId(req.Form.Get("id"))  // when creating a page, id == ""
	title := req.Form.Get("title")
	bodyFromEdit := req.Form.Get("body")
	body := server.syntaxHandler.EditToBody(bodyFromEdit)
	page := &Page{Id: id, Title: title, Body: body}

	if id == "" {
		id, err = server.pageStore.Create(page)
	} else {
		err = server.pageStore.Update(page)
	}

	if err != nil {
		internalError(res, err)
		return
	}

	http.Redirect(res, req, VIEW_ENTRYPOINT_PATH+string(id), http.StatusFound)
}

func (server *Server) handleDelete(res http.ResponseWriter, req *http.Request) {
	id, err := getRequestedPageId(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = server.pageStore.Delete(id)
	if err != nil {
		internalError(res, err)
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

func internalError(res http.ResponseWriter, err error) {
	http.Error(res, "internal error", http.StatusInternalServerError)  // do not return internal error details to client
	log.Println(err)
}

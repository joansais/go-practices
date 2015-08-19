package wikix

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
	. "github.com/joansais/go-practices/exception"
)

const (
	LIST_ENTRYPOINT_PATH   = "/"
	VIEW_ENTRYPOINT_PATH   = "/view/"
	CREATE_ENTRYPOINT_PATH = "/create/"
	EDIT_ENTRYPOINT_PATH   = "/edit/"
	SAVE_ENTRYPOINT_PATH   = "/save/"
	DELETE_ENTRYPOINT_PATH = "/delete/"
	HTML_TEMPLATE_FILES    = "/html/*.tmpl"
)

var (
	pageRequestPattern = regexp.MustCompile(`^/(view|edit|delete)/([a-zA-Z0-9]+)$`)
)

type Server struct {
	pageStore     PageStore
	syntaxHandler SyntaxHandler
	htmlTemplates *template.Template
}

func NewServer(store PageStore, syntax SyntaxHandler, assetsDir string) *Server {
	return &Server{
		pageStore:     store,
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
	defer Catch(errorHandler(res))

	ids := server.pageStore.ListAll()

	pageList := make(PageListModel, len(ids))
	for k, id := range ids {
		page := server.pageStore.Read(id)
		pageList[k] = &PageModel{Id: id, Title: page.Title}
	}
	sort.Sort(pageList)

	ThrowIf(server.htmlTemplates.ExecuteTemplate(res, "list", pageList))
}

func (server *Server) handleView(res http.ResponseWriter, req *http.Request) {
	defer Catch(errorHandler(res))

	page := server.pageStore.Read(getRequestedPageId(req))
	bodyAsHtml := server.syntaxHandler.BodyToHtml(page.Body)

	pageModel := &PageModel{Id: page.Id, Title: page.Title, BodyAsHtml: bodyAsHtml}
	ThrowIf(server.htmlTemplates.ExecuteTemplate(res, "view", pageModel))
}

func (server *Server) handleCreate(res http.ResponseWriter, req *http.Request) {
	defer Catch(errorHandler(res))

	title := req.URL.Query().Get("title")
	if title == "" {
		title = "New Page"
	}

	pageModel := &PageModel{Title: title}
	ThrowIf(server.htmlTemplates.ExecuteTemplate(res, "create", pageModel))
}

func (server *Server) handleEdit(res http.ResponseWriter, req *http.Request) {
	defer Catch(errorHandler(res))

	page := server.pageStore.Read(getRequestedPageId(req))
	bodyToEdit := server.syntaxHandler.BodyToEdit(page.Body)

	pageModel := &PageModel{Id: page.Id, Title: page.Title, BodyToEdit: bodyToEdit}
	ThrowIf(server.htmlTemplates.ExecuteTemplate(res, "edit", pageModel))
}

func (server *Server) handleSave(res http.ResponseWriter, req *http.Request) {
	defer Catch(errorHandler(res))

	err := req.ParseForm()
	if err != nil {
		Throw(InvalidRequestError{err})
	}

	id := PageId(req.Form.Get("id")) // when creating a page, id == ""
	title := req.Form.Get("title")
	bodyFromEdit := req.Form.Get("body")
	body := server.syntaxHandler.EditToBody(bodyFromEdit)
	page := &Page{Id: id, Title: title, Body: body}

	if id == "" {
		id = server.pageStore.Create(page)
	} else {
		server.pageStore.Update(page)
	}

	http.Redirect(res, req, VIEW_ENTRYPOINT_PATH+string(id), http.StatusFound)
}

func (server *Server) handleDelete(res http.ResponseWriter, req *http.Request) {
	defer Catch(errorHandler(res))
	server.pageStore.Delete(getRequestedPageId(req))
	http.Redirect(res, req, LIST_ENTRYPOINT_PATH, http.StatusFound)
}

func getRequestedPageId(req *http.Request) PageId {
	submatches := pageRequestPattern.FindStringSubmatch(req.URL.Path)
	if submatches == nil {
		Throw(InvalidRequestError{errors.New("invalid page id")})
	}
	return PageId(submatches[2])
}

type InvalidRequestError struct {
	cause error
}

func (err InvalidRequestError) Error() string {
	if err.cause != nil {
		return "invalid request: " + err.cause.Error()
	} else {
		return "invalid request"
	}
}

func errorHandler(res http.ResponseWriter) Catcher {
	return func(err error) {
		switch err.(type) {
		case InvalidRequestError:
			http.Error(res, err.Error(), http.StatusBadRequest)
		case UnexistentPageError:
			http.Error(res, err.Error(), http.StatusNotFound)
		default:
			http.Error(res, "internal error", http.StatusInternalServerError) // do not return internal error details to client
			log.Println(err)
		}
	}
}

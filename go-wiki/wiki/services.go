package wiki

import (
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

const (
	LIST_ENTRYPOINT_PATH = "/"
	VIEW_ENTRYPOINT_PATH = "/view/"
	EDIT_ENTRYPOINT_PATH = "/edit/"
	SAVE_ENTRYPOINT_PATH = "/save/"
	REMOVE_ENTRYPOINT_PATH = "/remove/"
	HTML_TEMPLATE_FILES  = "/html/*.tmpl"
)

var (
	actionPattern = regexp.MustCompile(`^/(view|edit|save|delete)/([a-zA-Z0-9]+)$`)
	htmlTemplates *template.Template
)

func SetAssetsDir(path string) {
	htmlTemplates = template.Must(template.ParseGlob(path + HTML_TEMPLATE_FILES))
}

func handleList(res http.ResponseWriter, req *http.Request) {
	titles, err := listPages()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTemplates.ExecuteTemplate(res, "list", titles)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleView(res http.ResponseWriter, req *http.Request) {
	title, err := getRequestedPageTitle(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := loadPage(title)
	if err != nil {
		http.Redirect(res, req, EDIT_ENTRYPOINT_PATH+title, http.StatusFound)
		return
	}

	err = htmlTemplates.ExecuteTemplate(res, "view", page)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEdit(res http.ResponseWriter, req *http.Request) {
	title, err := getRequestedPageTitle(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	err = htmlTemplates.ExecuteTemplate(res, "edit", page)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSave(res http.ResponseWriter, req *http.Request) {
	title, err := getRequestedPageTitle(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	body := req.Form.Get("body")
	page := &Page{Title: title, Body: body}
	err = page.save()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, VIEW_ENTRYPOINT_PATH+page.Title, http.StatusFound)
}

func handleRemove(res http.ResponseWriter, req *http.Request) {
	title, err := getRequestedPageTitle(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := loadPage(title)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	err = page.remove()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, LIST_ENTRYPOINT_PATH, http.StatusFound)
}

func getRequestedPageTitle(req *http.Request) (string, error) {
	submatches := actionPattern.FindStringSubmatch(req.URL.Path)
	if submatches == nil {
		return "", errors.New("Invalid page title")
	}
	return submatches[2], nil
}

func RegisterServices() {
	http.HandleFunc(LIST_ENTRYPOINT_PATH, handleList)
	http.HandleFunc(VIEW_ENTRYPOINT_PATH, handleView)
	http.HandleFunc(EDIT_ENTRYPOINT_PATH, handleEdit)
	http.HandleFunc(SAVE_ENTRYPOINT_PATH, handleSave)
	http.HandleFunc(REMOVE_ENTRYPOINT_PATH, handleRemove)
}

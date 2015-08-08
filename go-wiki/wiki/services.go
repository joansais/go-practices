package wiki

import (
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

const (
	VIEW_ENTRYPOINT_PATH = "/view/"
	EDIT_ENTRYPOINT_PATH = "/edit/"
	SAVE_ENTRYPOINT_PATH = "/save/"
	HTML_TEMPLATE_FILES  = "assets/views/*.html.template"
)

var (
	entrypointPattern = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	htmlTemplates = template.Must(template.ParseGlob(HTML_TEMPLATE_FILES))
)

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

	body := template.HTML(req.Form.Get("body"))
	page := &Page{Title: title, Body: body}
	err = page.save()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, VIEW_ENTRYPOINT_PATH+page.Title, http.StatusFound)
}

func getRequestedPageTitle(req *http.Request) (string, error) {
	submatches := entrypointPattern.FindStringSubmatch(req.URL.Path)
	if submatches == nil {
		return "", errors.New("Invalid page title")
	}
	return submatches[2], nil
}

func RegisterServices() {
	http.HandleFunc(VIEW_ENTRYPOINT_PATH, handleView)
	http.HandleFunc(EDIT_ENTRYPOINT_PATH, handleEdit)
	http.HandleFunc(SAVE_ENTRYPOINT_PATH, handleSave)
}

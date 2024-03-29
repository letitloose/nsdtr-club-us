package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		IsActive:        app.isActive(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	var ts *template.Template
	var ok bool
	var err error

	if !app.useTemplateCache {
		ts, err = getTemplateSet(page)
		if err != nil {
			app.serverError(w, err)
			err := fmt.Errorf("the template %s does not exist", page)
			app.serverError(w, err)
			return
		}
	} else {
		app.infoLog.Println("using template cache.")
		ts, ok = app.templateCache[page]
		if !ok {
			err := fmt.Errorf("the template %s does not exist", page)
			app.serverError(w, err)
			return
		}
	}

	buf := new(bytes.Buffer)
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) isActive(r *http.Request) bool {
	isActive, ok := r.Context().Value(isActiveContextKey).(bool)
	if !ok {
		return false
	}

	return isActive
}

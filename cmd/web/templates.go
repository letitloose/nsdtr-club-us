package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/letitloose/nsdtr-club-us/internal/models"
	"github.com/letitloose/nsdtr-club-us/ui"
)

type templateData struct {
	CurrentYear     int
	Member          *models.Member
	Members         []*models.Member
	Form            any
	Flash           string
	IsAuthenticated bool
	IsActive        bool
	CSRFToken       string
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("01/02/2006")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func getTemplateSet(page string) (*template.Template, error) {

	// Create a slice containing the filepath patterns for the templates we
	// want to parse.
	patterns := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/" + page,
	}

	// Use ParseFS() instead of ParseFiles() to parse the template files
	// from the ui.Files embedded filesystem.
	ts, err := template.New(page).Funcs(functions).ParseFiles(patterns...)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

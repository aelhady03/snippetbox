package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/aelhady03/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
}

type templateCache map[string]*template.Template

var functions = template.FuncMap{
	"humanDate": humanDate,
	"hello":     func() string { return "hello, world" },
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func newTemplateCache() (templateCache, error) {
	cache := make(templateCache)

	base, err := template.New("").Funcs(functions).ParseFiles("./ui/html/base.tmpl")
	if err != nil {
		return nil, err
	}

	// associates the resulting templates with base
	if _, err = base.ParseGlob("./ui/html/partials/*.tmpl"); err != nil {
		return nil, err
	}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		// clone the base so each page gets a fresh template set.
		ts, err := base.Clone()
		if err != nil {
			return nil, err
		}

		// re-use the same base while parsing the page
		// then for each page cache their own template set.
		if _, err = ts.ParseFiles(page); err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}

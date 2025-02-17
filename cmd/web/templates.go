package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"PhoceeneAuto/ui"
)

var functions = template.FuncMap{
	"humanDate":     humanDate,
	"bytesToString": bytesToString,
	"increment":     increment,
	"decrement":     decrement,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func bytesToString(b []byte) string {
	if b != nil {
		return string(b)
	}
	return ""
}

func increment(n int) int {
	return n + 1
}

func decrement(n int) int {
	return n - 1
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "templates/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"templates/base.tmpl",
			"templates/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

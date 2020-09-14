package main

import (
	"html/template"
	"ntguilty.me/notes-app/pkg/models"
	"path/filepath"
	"time"
)

//Define a templateData type to act as the holding structure for
//any dynamic data that we want to pass to our HTML templates (because our htmpl template can

type templateData struct {
	CurrentYear int
	Note  *models.Note
	Notes []*models.Note
}


// Custom template functions can accept as many parameters as they need to,
// but they must return one value only. The only exception to this is
// if you want to return an error as the second value.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate" : humanDate,

}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	//Get a slice of all filepaths with the extension to gives us a slice of all
	//the 'page' templates for the app.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' or 'partial' templates to the
		// template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}

	return cache, nil
}

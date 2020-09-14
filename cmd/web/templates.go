package main

import (
	"html/template"
	"ntguilty.me/notes-app/pkg/models"
	"path/filepath"
)

//Define a templateData type to act as the holding structure for
//any dynamic data that we want to pass to our HTML templates (because our htmpl template can

type templateData struct {
	CurrentYear int
	Note  *models.Note
	Notes []*models.Note
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

		ts, err := template.ParseFiles(page)
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

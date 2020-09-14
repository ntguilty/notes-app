package main

import "ntguilty.me/notes-app/pkg/models"

//Define a templateData type to act as the holding structure for
//any dynamic data that we want to pass to our HTML templates (because our htmpl template can

type templateData struct {
	Note  *models.Note
	Notes []*models.Note
}

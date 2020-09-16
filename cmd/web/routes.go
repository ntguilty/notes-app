package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/note", app.showNote)
	mux.HandleFunc("/note/create", app.createNote)
	//Create a fileserver to serve files from "./ui/static/", register a file server as a handle for "/static/" and strip to match paths.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return app.logRequest(secureHeaders(mux))
}

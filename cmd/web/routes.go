package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showNote)
	mux.HandleFunc("/snippet/create", app.createNote)
	//Create a fileserver to serve files from "./ui/static/", register a file server as a handle for "/static/" and strip to match paths.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}

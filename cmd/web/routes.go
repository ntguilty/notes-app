package main

import (
	"github.com/justinas/alice"
	"net/http"
	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	basicMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/note/create", dynamicMiddleware.ThenFunc(app.createNoteForm))
	mux.Post("/note/create", dynamicMiddleware.ThenFunc(app.createNote))
	// Moved down because /snippet/create also match /snipet/:id and Pat
	// package matches pattern in the order that they are registered.
	mux.Get("/note/:id", dynamicMiddleware.ThenFunc(app.showNote))

	//Create a fileserver to serve files from "./ui/static/", register a file server as a handle for "/static/" and strip to match paths.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return basicMiddleware.Then(mux)
}

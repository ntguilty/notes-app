package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	basicMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/note/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNoteForm))
	mux.Post("/note/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNote))
	// Moved down because /snippet/create also match /snipet/:id and Pat
	// package matches pattern in the order that they are registered.
	mux.Get("/note/:id", dynamicMiddleware.ThenFunc(app.showNote))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	//Create a fileserver to serve files from "./ui/static/", register a file server as a handle for "/static/" and strip to match paths.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return basicMiddleware.Then(mux)
}

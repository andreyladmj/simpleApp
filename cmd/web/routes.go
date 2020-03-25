package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.session.Enable)
	dynamicMiddleware := alice.New(app.authenticate)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.home))
	mux.Get("/checkAuth", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.checkAuth))
	mux.Get(os.Getenv("RETURN_URL"), http.HandlerFunc(app.returnUrl))
	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

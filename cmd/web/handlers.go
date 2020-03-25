package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home Page"))
}

func (app *application) checkAuth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("checkAuth Page"))
}

func (app *application) returnUrl(w http.ResponseWriter, r *http.Request) {
	token, ok := r.URL.Query()["token"]
	app.session.Put(r, "already_got_auth_answer", true)

	if !ok {
		app.errorLog.Printf("Auth by returned toked failed: %v", token)
	}

	app.session.Put(r, "token", token[0])

	fmt.Println("app.session.Exists)", app.session.Exists(r,"prevPage"))
	fmt.Println("app.session.GET)", app.session.GetString(r, "prevPage"))

	if !app.session.Exists(r,"prevPage") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	prevPage := app.session.GetString(r, "prevPage")

	http.Redirect(w, r, prevPage, http.StatusFound)
}

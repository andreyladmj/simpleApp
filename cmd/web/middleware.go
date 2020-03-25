package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
)


func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, "token")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		token := app.session.GetString(r, "token")
		fmt.Println("Got token", token)

		user := app.grpcClient.GetUser(token)
		fmt.Println("Got user", user)

		if user == nil {
			app.session.Remove(r, "token")
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			exist := app.session.Exists(r, "already_got_auth_answer")
			if exist {
				app.session.Remove(r, "already_got_auth_answer")
				app.serverError(w, errors.New("auth service already responded"))
				return
			}

			scheme := "https"

			if r.TLS == nil {
				scheme = "http"
			}

			fmt.Println("Save Page", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path))
			app.session.Put(r, "prevPage", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path))
			authService := fmt.Sprintf("%s?returnUrl=%s", os.Getenv("AUTH_SERVICE"), fmt.Sprintf("%s://%s%s", scheme, r.Host, os.Getenv("RETURN_URL")))
			http.Redirect(w, r, authService, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

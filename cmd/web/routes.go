package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// mux := http.NewServeMux()

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static")))
	// fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// unprotected routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/activate", dynamic.ThenFunc(app.activateUser))

	//protected routes
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	//active routes
	active := dynamic.Append(app.requireActive)
	router.Handler(http.MethodGet, "/member/create", active.ThenFunc(app.memberForm))
	router.Handler(http.MethodPost, "/member/create", active.ThenFunc(app.memberCreate))
	router.Handler(http.MethodGet, "/member/view/:id", active.ThenFunc(app.memberView))
	router.Handler(http.MethodGet, "/member", active.ThenFunc(app.memberList))
	router.Handler(http.MethodGet, "/membership/create", active.ThenFunc(app.membershipForm))
	router.Handler(http.MethodPost, "/membership/create", active.ThenFunc(app.membershipCreate))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

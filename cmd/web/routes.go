package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/letitloose/nsdtr-club-us/ui"
)

func (app *application) routes() http.Handler {
	// mux := http.NewServeMux()

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// unprotected routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/activate", dynamic.ThenFunc(app.userActivate))

	//protected routes
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/member/create", protected.ThenFunc(app.memberForm))
	router.Handler(http.MethodPost, "/member/create", protected.ThenFunc(app.memberCreate))
	router.Handler(http.MethodGet, "/member/view/:id", protected.ThenFunc(app.memberView))
	router.Handler(http.MethodGet, "/member", protected.ThenFunc(app.memberList))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

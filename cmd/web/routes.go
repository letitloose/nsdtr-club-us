package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/member/create", app.memberCreate)
	mux.HandleFunc("/member/view", app.memberView)
	mux.HandleFunc("/member", app.memberList)
	return mux
}

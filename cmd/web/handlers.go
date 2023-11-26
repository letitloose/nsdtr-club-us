package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/letitloose/nsdtr-club-us/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) memberForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create new member form"))
}

func (app *application) memberCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	firstname := "Lou"
	lastname := "Garwood"
	phone := "518-495-2003"
	email := "louis.garwood@gmail.com"
	website := "www.github.com/letitloose"
	region := 1
	joined := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)

	id, err := app.members.Insert(firstname, lastname, phone, email, website, region, joined)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/member/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) memberView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	member, err := app.members.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Member = member

	app.render(w, http.StatusOK, "member-view.html", data)
}

func (app *application) memberList(w http.ResponseWriter, r *http.Request) {

	members, err := app.members.List()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Members = members

	app.render(w, http.StatusOK, "member-list.html", data)

}

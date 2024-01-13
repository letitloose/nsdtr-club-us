package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/letitloose/nsdtr-club-us/internal/models"
	"github.com/letitloose/nsdtr-club-us/internal/services"
	"github.com/letitloose/nsdtr-club-us/internal/validator"
)

type memberCreateForm struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Website     string
	Region      int
	JoinedDate  string
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) memberForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = memberCreateForm{
		JoinedDate: time.Now().Format("2006-01-02"),
	}

	app.render(w, http.StatusOK, "member-create.html", data)
}

func (app *application) memberCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	region, err := strconv.Atoi(r.PostForm.Get("region"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := memberCreateForm{
		FirstName:   r.PostForm.Get("firstname"),
		LastName:    r.PostForm.Get("lastname"),
		PhoneNumber: r.PostForm.Get("phonenumber"),
		Email:       r.PostForm.Get("email"),
		Website:     r.PostForm.Get("website"),
		Region:      region,
		JoinedDate:  r.PostForm.Get("joindate"),
	}

	//validate
	form.CheckField(validator.NotBlank(form.FirstName), "firstname", "You must enter a first name.")
	form.CheckField(validator.NotBlank(form.LastName), "lastname", "You must enter a last name.")
	form.CheckField(validator.ValidEmail(form.Email), "email", "You must enter a valid email: name@domain.ext")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "member-create.html", data)
		return
	}

	joined, err := time.Parse("2006-01-02", r.PostForm.Get("joindate"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.members.Insert(form.FirstName, form.LastName, form.PhoneNumber, form.Email, form.Website, form.Region, joined)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//put success message in flash
	app.sessionManager.Put(r.Context(), "flash", "Member successfully created!")

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/member/view/%d", id), http.StatusSeeOther)
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

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = services.UserForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := services.UserForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	err = app.userService.InsertUser(&form)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else if errors.Is(err, models.ErrBadData) {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			return
		} else {
			app.serverError(w, err)
		}

		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please check your email to activate your account.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = services.UserForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := services.UserForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	id, err := app.userService.AuthenticateUser(&form)
	if err != nil {
		if errors.Is(err, models.ErrBadData) {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return
		} else if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) activateUser(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		app.notFound(w)
	}

	err := app.userService.ActivateUser(hash)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "activated.html", data)
}

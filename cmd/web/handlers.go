package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/letitloose/nsdtr-club-us/internal/models"
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

type userSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userSignupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidEmail(form.Email), "email", "You must enter a valid email: name@domain.ext")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	//user created successfully,  send an email with the validation link
	verificationHash, err := app.users.GetVerificationHashByEmail(form.Email)
	body := fmt.Sprintf(
		`<html>
			<body>
				<h1>Hello %s!</h1>
				<p>Please <a href="https://localhost:8080/user/activate?hash=%s">click here</a> to validate your email and activate your account.<p>
			</body>
		</html>`, form.Name, verificationHash)
	err = app.email.SendEmail("Welcome to NSDTRC-USA Membership", "", body)
	if err != nil {
		app.serverError(w, err)
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Please enter your email to login")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
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

func (app *application) userActivate(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		app.notFound(w)
	}

	userID, err := app.users.GetByVerificationHash(hash)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.users.Activate(userID)
	if err != nil {
		app.serverError(w, err)
	}

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "activated.html", data)
}

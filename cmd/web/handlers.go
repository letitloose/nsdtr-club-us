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
	data := app.newTemplateData(r)

	data.Form = services.MemberForm{
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

	form := &services.MemberForm{
		FirstName:   r.PostForm.Get("firstname"),
		LastName:    r.PostForm.Get("lastname"),
		PhoneNumber: r.PostForm.Get("phonenumber"),
		Email:       r.PostForm.Get("email"),
		Website:     r.PostForm.Get("website"),
		Region:      region,
		JoinedDate:  r.PostForm.Get("joindate"),
	}

	id, err := app.memberService.CreateMember(form)
	if err != nil {
		if errors.Is(err, models.ErrBadData) {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "member-create.html", data)
			return
		} else {
			app.serverError(w, err)
		}
	}

	//put success message in flash
	app.sessionManager.Put(r.Context(), "flash", "Member successfully created!")

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/member/view/%d", id), http.StatusSeeOther)
}

func (app *application) membershipForm(w http.ResponseWriter, r *http.Request) {
	// params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(r.URL.Query().Get("memberID"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	member, err := app.memberService.GetMemberProfile(id)
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

	data.Form = services.MembershipForm{
		MemberID: id,
		Year:     2024,
	}

	app.render(w, http.StatusOK, "membership-create.html", data)
}

func (app *application) membershipCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	memberID, err := strconv.Atoi(r.PostForm.Get("member-id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(r.PostForm.Get("year"))
	if err != nil {
		app.errorLog.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	membershipAmount, err := strconv.Atoi(r.PostForm.Get("membership-amount"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	rosterAmt, err := strconv.Atoi(r.PostForm.Get("roster-amount"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	healthDonation, err := strconv.Atoi(r.PostForm.Get("health-donation"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	rescueDonation, err := strconv.Atoi(r.PostForm.Get("rescue-donation"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &services.MembershipForm{
		MemberID:         memberID,
		Year:             year,
		MembershipType:   r.PostForm.Get("membership-type"),
		MembershipAmount: float32(membershipAmount),
		RosterAmount:     float32(rosterAmt),
		HealthDonations:  float32(healthDonation),
		RescueDonations:  float32(rescueDonation),
	}

	err = app.memberService.AddMembership(form)
	if err != nil {
		if errors.Is(err, models.ErrBadData) {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "membership-create.html", data)
			return
		} else {
			app.serverError(w, err)
		}
	}

	//put success message in flash
	app.sessionManager.Put(r.Context(), "flash", "Membership successfully created!")

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/member/view/%d", memberID), http.StatusSeeOther)
}

func (app *application) memberView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	member, err := app.memberService.GetMemberProfile(id)
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

	members, err := app.memberService.List()
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

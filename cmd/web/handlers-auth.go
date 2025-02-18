package main

import (
	"errors"
	"fmt"
	"net/http"

	"PhoceeneAuto/internal/data"
	"PhoceeneAuto/internal/validator"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Login"

	// filling the form with empty values
	tmplData.Form = newUserLoginForm()

	// rendering the template
	app.render(w, r, http.StatusOK, "login.tmpl", tmplData)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newUserLoginForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	form.Check(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.Check(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	if form.ValidatePassword(form.Password); !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "login.tmpl")
		return
	}

	// fetching the user with the mail address
	user, err := app.models.UserModel.GetByEmail(form.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			form.AddNonFieldError("invalid credentials")
			app.failedValidationError(w, r, form, &form.Validator, "login.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// matching the password
	match, err := user.Password.Matches(form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking the password match
	if !match {
		form.AddNonFieldError("invalid credentials")
		app.failedValidationError(w, r, form, &form.Validator, "login.tmpl")
		return
	}

	// renewing the user session
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// storing the user id and user role in the user session
	app.sessionManager.Put(r.Context(), authenticatedUserIDSessionManager, user.ID)
	app.sessionManager.Put(r.Context(), userRoleSessionManager, user.Role)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {

	// logging the user out
	err := app.logout(r)
	if err != nil {

		// DEBUG
		app.logger.Debug(fmt.Sprintf("error: %s", err.Error()))

		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

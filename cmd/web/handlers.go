package main

import (
	"PhoceeneAuto/internal/data"
	"PhoceeneAuto/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

/* #############################################################################
/*	COMMON
/* #############################################################################*/

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Not Found"

	// rendering the template
	app.render(w, r, http.StatusOK, "error.tmpl", tmplData)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Oooops"

	// setting the error title and message
	tmplData.Error.Title = "Error 405"
	tmplData.Error.Message = "Something went wrong!"

	// rendering the template
	app.render(w, r, http.StatusOK, "error.tmpl", tmplData)
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {

	// checking the query
	if r.URL.Query() == nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Search"

	// retrieving the research text
	tmplData.Search = r.URL.Query().Get("q")

	// search in the posts
	var err error
	tmplData.Posts.List, tmplData.Posts.Metadata, err = app.models.PostModel.Get(tmplData.Search, data.NewPostFilters(r.URL.Query()))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// rendering the template
	app.render(w, r, http.StatusOK, "search.tmpl", tmplData)
}

/* #############################################################################
/*	USER ACCESS
/* #############################################################################*/

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Register"

	// filling the form with empty values
	tmplData.Form = newUserCreateForm()

	// rendering the template
	app.render(w, r, http.StatusOK, "create-user.tmpl", tmplData)
}

func (app *application) createUserPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newUserCreateForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// DEBUG
	app.logger.Debug(fmt.Sprintf("form: %+v", form))

	// checking the data from the user
	form.StringCheck(form.Username, 2, 70, true, "username")
	form.ValidateEmail(form.Email)
	form.ValidateRegisterPassword(form.Password, form.ConfirmPassword)

	// return to create user page if there is an error
	if !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "create-user.tmpl")
		return
	}

	// creating the user
	user := &data.User{
		Name:   form.Username,
		Email:  form.Email,
		Status: data.UserToActivate,
	}

	// setting the password hash
	err = user.Password.Set(form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// verifying the user data
	if data.ValidateUser(&form.Validator, user); !form.Valid() {

		// redirect to login page with errors
		app.failedValidationError(w, r, form, &form.Validator, "create-user.tmpl")
		return
	}

	// inserting the user in the DB
	err = app.models.UserModel.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			form.AddFieldError("email", "a user with this email address already exists")
			app.failedValidationError(w, r, form, &form.Validator, "create-user.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// Generating an activation token to send it via mail to the user
	token, err := app.models.TokenModel.New(user.ID, 3*24*time.Hour, data.TokenActivation)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.background(func() {

		mailData := map[string]any{
			"userID":          user.ID,
			"username":        user.Name,
			"activationToken": token.Plaintext,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", mailData)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	app.sessionManager.Put(r.Context(), "flash", "We've sent you a confirmation email!")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

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

	// storing the user id in the user session
	app.sessionManager.Put(r.Context(), authenticatedUserIDSessionManager, user.ID)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

/* #############################################################################
/*	RESTRICTED
/* #############################################################################*/

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Dashboard"

	// rendering the template
	app.render(w, r, http.StatusOK, "dashboard.tmpl", tmplData)
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

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Update user"

	// retrieving user ID
	id := app.getUserID(r)

	// fetching user data
	user, err := app.models.UserModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// filling the form with user values
	tmplData.Form = newUserUpdateForm(user)

	// rendering the template
	app.render(w, r, http.StatusOK, "user-update.tmpl", tmplData)
}

func (app *application) updateUserPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newUserUpdateForm(nil)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// fetching the authenticated user
	user, err := app.models.UserModel.GetByID(app.getUserID(r))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking the data from the user
	var isEmpty = true
	if form.Username != nil {
		isEmpty = false
		form.StringCheck(*form.Username, 2, 70, false, "username")
		user.Name = *form.Username
	}
	if form.Password != nil || form.NewPassword != nil || form.ConfirmationPassword != nil {
		isEmpty = false
		form.ValidateNewPassword(*form.NewPassword, *form.ConfirmationPassword)
		err = user.Password.Set(*form.NewPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	if form.Email != nil {
		isEmpty = false
		form.ValidateEmail(*form.Email)
		user.Email = *form.Email
	}
	if isEmpty {
		form.AddNonFieldError("at least one field is required")
	}

	// return to update-user page if there is an error
	if !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "user-update.tmpl")
		return
	}

	// updating the user
	err = app.models.UserModel.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		case errors.Is(err, data.ErrDuplicateEmail):
			form.AddFieldError("email", "email is already in use")
			app.failedValidationError(w, r, form, &form.Validator, "user-update.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your data has been updated successfully!")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

/* #############################################################################
/*	AJAX CALLS
/* #############################################################################*/

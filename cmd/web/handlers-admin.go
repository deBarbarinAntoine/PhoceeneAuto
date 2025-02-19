package main

import (
	"errors"
	"fmt"
	"net/http"

	"PhoceeneAuto/internal/data"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Create User"

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

	// creating the user with the form data
	user := form.toUser()

	// setting the password hash
	err = user.Password.Set(form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking the password
	form.ValidateRegisterPassword(form.Password, form.ConfirmPassword)

	// return to create user page if there is an error
	if !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "create-user.tmpl")
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

	// adding notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("User %s has been created successfully", user.Name))

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Update user"

	// retrieving ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

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

	// getting the user id
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// checking that the path id is equal to the form id
	if id != *form.ID {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// fetching the user to update
	user, err := app.models.UserModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking the data from the user
	isEmpty := form.toUser(user)
	if form.Password != nil || form.NewPassword != nil || form.ConfirmationPassword != nil {
		isEmpty = false
		form.ValidateNewPassword(*form.NewPassword, *form.ConfirmationPassword)
		err = user.Password.Set(*form.NewPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	if isEmpty {
		form.AddNonFieldError("at least one field is required")
	}
	data.ValidateUser(&form.Validator, user)

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

	// check if updated user is the one logged in
	if app.getUserID(r) == user.ID {

		// update user role in the SessionManager
		app.sessionManager.Put(r.Context(), userRoleSessionManager, user.Role)
	}

	app.sessionManager.Put(r.Context(), "flash", "Your data has been updated successfully!")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) userGet(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) reports(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

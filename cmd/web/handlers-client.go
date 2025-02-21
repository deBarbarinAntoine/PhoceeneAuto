package main

import (
	"errors"
	"fmt"
	"net/http"
	
	"PhoceeneAuto/internal/data"
)

func (app *application) clientGet(w http.ResponseWriter, r *http.Request) {
	
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Client"
	
	// retrieving the Client ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	
	// retrieving the Client
	client, err := app.models.ClientModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	
	// adding the client in the template data
	tmplData.Client = *client
	
	// rendering the template
	app.render(w, r, http.StatusOK, "client.tmpl", tmplData)
}

func (app *application) deleteClient(w http.ResponseWriter, r *http.Request) {
	
	// retrieving the Client ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	
	// retrieving the Client
	client, err := app.models.ClientModel.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	
	// deleting the Client
	err = app.models.ClientModel.Delete(client)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	
	// adding the notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Client %s %s has been deleted successfully!", client.FirstName, client.LastName))
	
	// redirecting to the dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) createClient(w http.ResponseWriter, r *http.Request) {
	
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Create Client"
	
	// filling the form with empty values
	tmplData.Form = newClientCreateForm()
	
	// rendering the template
	app.render(w, r, http.StatusOK, "client-create.tmpl", tmplData)
}

func (app *application) createClientPost(w http.ResponseWriter, r *http.Request) {
	
	// retrieving the form data
	form := newClientCreateForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	
	// DEBUG
	app.logger.Debug(fmt.Sprintf("form: %+v", form))
	
	// creating the client
	client := form.toClient()
	
	// verifying the client data
	if data.ValidateClient(&form.Validator, client); !form.Valid() {
		
		// redirect to form with errors
		app.failedValidationError(w, r, form, &form.Validator, "client-create.tmpl")
		return
	}
	
	// inserting the client in the DB
	err = app.models.ClientModel.Insert(client)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			form.AddFieldError("email", "a client with this email address already exists")
			app.failedValidationError(w, r, form, &form.Validator, "client-create.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}
	
	// adding a notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Client %s %s has been created successfully!", client.FirstName, client.LastName))
	
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) updateClient(w http.ResponseWriter, r *http.Request) {
	
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Update Client"
	
	// retrieving ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	
	// fetching client data
	client, err := app.models.ClientModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	
	// filling the form with client values
	tmplData.Form = newClientUpdateForm(client)
	
	// rendering the template
	app.render(w, r, http.StatusOK, "client-update.tmpl", tmplData)
}

func (app *application) updateClientPost(w http.ResponseWriter, r *http.Request) {
	
	// retrieving the ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	
	// retrieving the Client
	client, err := app.models.ClientModel.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	
	// retrieving the form data
	form := newClientUpdateForm(client)
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	
	// DEBUG
	app.logger.Debug(fmt.Sprintf("form: %+v", form))
	
	// creating the client
	form.toClient(client)
	
	// verifying the client data
	if data.ValidateClient(&form.Validator, client); !form.Valid() {
		
		// redirect to form with errors
		app.failedValidationError(w, r, form, &form.Validator, "client-create.tmpl")
		return
	}
	
	// updating the client in the DB
	err = app.models.ClientModel.Update(client)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			form.AddFieldError("email", "a client with this email address already exists")
			app.failedValidationError(w, r, form, &form.Validator, "client-create.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}
	
	// adding a notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Client %s %s has been updated successfully!", client.FirstName, client.LastName))
	
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

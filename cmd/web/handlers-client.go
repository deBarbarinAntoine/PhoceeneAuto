package main

import "net/http"

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

	tmplData.Client = *client

	// rendering the template
	app.render(w, r, http.StatusOK, "create-user.tmpl", tmplData)
}

func (app *application) deleteClient(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) createClient(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Create Client"

	// filling the form with empty values
	tmplData.Form = newClientCreateForm()

	// rendering the template
	app.render(w, r, http.StatusOK, "create-client.tmpl", tmplData)
}

func (app *application) createClientPost(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) updateClient(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) updateClientPost(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

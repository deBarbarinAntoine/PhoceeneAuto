package main

import (
	"net/http"
)

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

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Dashboard"

	// rendering the template
	app.render(w, r, http.StatusOK, "dashboard.tmpl", tmplData)
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// Decode form values into searchForm struct
	var form searchForm
	decoder := form.NewDecoder()
	if err := decoder.Decode(&form, r.Form); err != nil {
		app.serverError(w, r, err)
		return
	}

	// Validate form inputs
	form.Validate()
	if !form.Valid() {
		app.clientError(w, r, http.StatusUnprocessableEntity)
		return
	}

	// Retrieve search results from DB (pseudo-query)
	results, err := app.models.Cars.Search(form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Prepare template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Search Results"

	// Render the template with results
	app.render(w, r, http.StatusOK, "search.tmpl", tmplData)
}

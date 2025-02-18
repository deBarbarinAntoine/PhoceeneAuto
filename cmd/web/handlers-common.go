package main

import (
	"errors"
	"net/http"

	"PhoceeneAuto/internal/data"
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

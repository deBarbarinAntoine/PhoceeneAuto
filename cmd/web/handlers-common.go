package main

import (
	"PhoceeneAuto/internal/data"
	"net/http"
	"strconv"
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

	var form data.Search

	form.Search = getString(r, "search")
	form.Make = getString(r, "make")
	form.Model = getString(r, "model")
	form.Year = getInt(r, "year")
	form.PriceMin = getFloat64(r, "price_min")
	form.PriceMax = getFloat64(r, "price_max")
	form.KmMin = getFloat64(r, "km_min")
	form.KmMax = getFloat64(r, "km_max")
	form.Color = getString(r, "color")
	form.Transmission = getString(r, "transmission")
	form.Fuel1 = getString(r, "fuel1")
	form.Fuel2 = getString(r, "fuel2")
	form.SizeClass = getString(r, "size_class")
	form.OwnerCount = getInt(r, "owner_count")
	form.Shop = getString(r, "shop")
	form.Status = getString(r, "status")
	form.ClientName = getString(r, "client_name")
	form.Email = getString(r, "email")
	form.Phone = getString(r, "phone")
	form.ClientStatus = getString(r, "client_status")
	form.TransactionID = getInt(r, "transaction_id")
	form.UserID = getInt(r, "user_id")
	form.TransactionStatus = getString(r, "transaction_status")
	form.DateStart = getString(r, "date_start")
	form.DateEnd = getString(r, "date_end")
	form.LeaseAmountMin = getFloat64(r, "lease_min")
	form.LeaseAmountMax = getFloat64(r, "lease_max")

	// retrieving the Client
	result, err := app.models.SearchModel.SearchAll(form)
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

// Helper functions to parse different types
func getString(r *http.Request, key string) *string {
	value := r.FormValue(key)
	if value == "" {
		return nil
	}
	return &value
}

func getInt(r *http.Request, key string) *int {
	value := r.FormValue(key)
	if value == "" {
		return nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &intValue
}

func getFloat64(r *http.Request, key string) *float64 {
	value := r.FormValue(key)
	if value == "" {
		return nil
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &floatValue
}

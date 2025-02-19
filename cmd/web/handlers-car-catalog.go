package main

import (
	"errors"
	"fmt"
	"net/http"

	"PhoceeneAuto/internal/data"
)

func (app *application) carCatalogGet(w http.ResponseWriter, r *http.Request) {
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Car Catalog"

	// retrieving the Car Catalog ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// fetching all cars in the catalog
	cars, err := app.models.CarCatalogModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// adding the cars to the template data
	tmplData.CarCatalog = *cars

	// rendering the template
	app.render(w, r, http.StatusOK, "car-catalog.tmpl", tmplData)
}

func (app *application) deleteCarCatalog(w http.ResponseWriter, r *http.Request) {
	// retrieving the Car ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// retrieving the Car Catalog
	car, err := app.models.CarCatalogModel.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// deleting the Car
	err = app.models.CarCatalogModel.Delete(car)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// adding the notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Car Catalog has been deleted successfully!"))

	// redirecting to the car catalog
	http.Redirect(w, r, "/car-catalog", http.StatusSeeOther)
}

func (app *application) createCarCatalog(w http.ResponseWriter, r *http.Request) {
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Create Car Catalog"

	// filling the form with empty values
	tmplData.Form = newCarCatalogCreateForm()

	// rendering the template
	app.render(w, r, http.StatusOK, "create-car-catalog.tmpl", tmplData)
}

func (app *application) createCarCatalogPost(w http.ResponseWriter, r *http.Request) {
	// retrieving the form data
	form := newCarCatalogCreateForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// DEBUG
	app.logger.Debug(fmt.Sprintf("form: %+v", form))

	// creating the user with the form data
	carCatalog := form.toCarCatalog()

	err = app.models.CarCatalogModel.Insert(carCatalog)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCarCatalog):
			form.Validator.AddFieldError("make", "This car catalog already exists")
			app.failedValidationError(w, r, form, &form.Validator, "create-car-catalog.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// adding notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Catalog has been created successfully"))

	http.Redirect(w, r, "/car_catalog", http.StatusSeeOther)
}

func (app *application) updateCarCatalog(w http.ResponseWriter, r *http.Request) {
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Update Car Catalog"

	// retrieving ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// fetching car data
	carCatalog, err := app.models.CarCatalogModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// filling the form with car values
	tmplData.Form = newCarCatalogUpdateForm(carCatalog)

	// rendering the template
	app.render(w, r, http.StatusOK, "car-update.tmpl", tmplData)
}

func (app *application) updateCarCatalogPost(w http.ResponseWriter, r *http.Request) {
	// retrieving the form data
	form := newCarCatalogUpdateForm(nil)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// getting the car catalog id
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

	// fetching the car to update
	carCatalog, err := app.models.CarCatalogModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking the data from the car
	isEmpty := form.toCarCatalog(carCatalog)
	if isEmpty {
		form.AddNonFieldError("at least one field is required")
	}

	// return to update car page if there is an error
	if !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "car-update.tmpl")
		return
	}

	// updating the car
	err = app.models.CarCatalogModel.Update(carCatalog)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		case errors.Is(err, data.ErrDuplicateCarCatalog):
			form.Validator.AddFieldError("make", "This car catalog already exists")
			app.failedValidationError(w, r, form, &form.Validator, "create-car-catalog.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// adding notification message
	app.sessionManager.Put(r.Context(), "flash", "Car data has been updated successfully!")

	http.Redirect(w, r, "/car-catalog", http.StatusSeeOther)
}

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
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Car Catalog %s has been deleted successfully!", car.Name))

	// redirecting to the car catalog
	http.Redirect(w, r, "/car-catalog", http.StatusSeeOther)
}

func (app *application) createCarCatalog(w http.ResponseWriter, r *http.Request) {

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
		case errors.Is(err, data.ErrDuplicateEmail):
			// TODO : check how to get the trigger error here
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

func (app *application) createCarCatalogPost(w http.ResponseWriter, r *http.Request) {
	// retrieving the form data
	form := newCarCreateForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// DEBUG
	app.logger.Debug(fmt.Sprintf("form: %+v", form))

	// creating the car with the form data
	car := form.toCar()

	// validating the car data
	form.ValidateCarDetails()

	// return to create car page if there is an error
	if !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "create-car.tmpl")
		return
	}

	// inserting the car into the DB
	err = app.models.CarCatalogModel.Insert(car)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCar):
			form.AddFieldError("name", "a car with this name already exists")
			app.failedValidationError(w, r, form, &form.Validator, "create-car.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// adding notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Car %s has been created successfully", car.Name))

	http.Redirect(w, r, "/car-catalog", http.StatusSeeOther)
}

func (app *application) updateCarCatalog(w http.ResponseWriter, r *http.Request) {
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Update Car"

	// retrieving ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// fetching car data
	car, err := app.models.CarCatalogModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// filling the form with car values
	tmplData.Form = newCarUpdateForm(car)

	// rendering the template
	app.render(w, r, http.StatusOK, "car-update.tmpl", tmplData)
}

func (app *application) updateCarCatalogPost(w http.ResponseWriter, r *http.Request) {
	// retrieving the form data
	form := newCarUpdateForm(nil)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// getting the car id
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
	car, err := app.models.CarCatalogModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking the data from the car
	isEmpty := form.toCar(car)
	if isEmpty {
		form.AddNonFieldError("at least one field is required")
	}

	// validating car data
	form.ValidateCarDetails()

	// return to update car page if there is an error
	if !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "car-update.tmpl")
		return
	}

	// updating the car
	err = app.models.CarCatalogModel.Update(car)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		case errors.Is(err, data.ErrDuplicateCar):
			form.AddFieldError("name", "car name is already in use")
			app.failedValidationError(w, r, form, &form.Validator, "car-update.tmpl")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// adding notification message
	app.sessionManager.Put(r.Context(), "flash", "Car data has been updated successfully!")

	http.Redirect(w, r, "/car-catalog", http.StatusSeeOther)
}

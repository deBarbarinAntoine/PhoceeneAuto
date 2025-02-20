package main

import (
	"errors"
	"fmt"
	"net/http"

	"PhoceeneAuto/internal/data"
)

func (app *application) carProductGet(w http.ResponseWriter, r *http.Request) {
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Car Product"

	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	carProduct, err := app.models.CarProductModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplData.CarProduct = *carProduct

	app.render(w, r, http.StatusOK, "car-product.tmpl", tmplData)
}

func (app *application) deleteCarProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	carProduct, err := app.models.CarProductModel.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.models.CarProductModel.Delete(carProduct)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Car Product %d has been deleted successfully!", carProduct.ID))
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) createCarProduct(w http.ResponseWriter, r *http.Request) {
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Create Car Product"
	tmplData.Form = newCarProductCreateForm()
	app.render(w, r, http.StatusOK, "create-car-product.tmpl", tmplData)
}

func (app *application) createCarProductPost(w http.ResponseWriter, r *http.Request) {
	form := newCarProductCreateForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	carProduct := form.toCarProduct()

	// check the Car Product data form
	if data.ValidateCarProduct(&form.Validator, *carProduct); !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "create-car-product.tmpl")
		return
	}

	err = app.models.CarProductModel.Insert(carProduct)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Car Product has been created successfully")
	http.Redirect(w, r, fmt.Sprintf("/car-product/%d", carProduct.ID), http.StatusSeeOther)
}

func (app *application) updateCarProduct(w http.ResponseWriter, r *http.Request) {
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Update Car Product"

	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	carProduct, err := app.models.CarProductModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplData.Form = newCarProductUpdateForm(carProduct)
	app.render(w, r, http.StatusOK, "car-update.tmpl", tmplData)
}

func (app *application) updateCarProductPost(w http.ResponseWriter, r *http.Request) {
	form := newCarProductUpdateForm(nil)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	if id != *form.ID {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	carProduct, err := app.models.CarProductModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	isEmpty := form.toCarProduct(carProduct)
	if isEmpty {
		form.AddNonFieldError("At least one field is required")
	}

	if data.ValidateCarProduct(&form.Validator, *carProduct); !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "car-update.tmpl")
		return
	}

	err = app.models.CarProductModel.Update(carProduct)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Car Product %d has been updated successfully!", carProduct.ID))
	http.Redirect(w, r, fmt.Sprintf("/car-product/%d", carProduct.ID), http.StatusSeeOther)
}

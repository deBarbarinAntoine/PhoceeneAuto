package main

import (
	"errors"
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

	app.sessionManager.Put(r.Context(), "flash", "Car Product has been deleted successfully!")
	http.Redirect(w, r, "/car-products", http.StatusSeeOther)
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

	err = app.models.CarProductModel.Insert(carProduct)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Product has been created successfully")
	http.Redirect(w, r, "/car-products", http.StatusSeeOther)
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

	if !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "car-update.tmpl")
		return
	}

	err = app.models.CarProductModel.Update(carProduct)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Car Product has been updated successfully!")
	http.Redirect(w, r, "/car-products", http.StatusSeeOther)
}

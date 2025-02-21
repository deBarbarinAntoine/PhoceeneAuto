package main

import (
	"errors"
	"fmt"
	"net/http"

	"PhoceeneAuto/internal/data"
)

func (app *application) transactionGet(w http.ResponseWriter, r *http.Request) {
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Transaction"

	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusNotFound)
		return
	}

	transaction, err := app.models.TransactionModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplData.Transaction = *transaction

	app.render(w, r, http.StatusOK, "transaction.tmpl", tmplData)

}

func (app *application) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	// retrieving the Car ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	transaction, err := app.models.TransactionModel.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.models.TransactionModel.Delete(transaction)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Transaction deleted successfully")

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func (app *application) createTransaction(w http.ResponseWriter, r *http.Request) {
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Transaction"

	tmplData.Form = newTransactionCreateForm()

	app.render(w, r, http.StatusOK, "transaction-create.tmpl", tmplData)
}

func (app *application) createTransactionPost(w http.ResponseWriter, r *http.Request) { // retrieving the form data
	form := newTransactionCreateForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	// DEBUG
	app.logger.Debug(fmt.Sprintf("form: %+v", form))

	// creating the user with the form data
	transaction := form.toTransaction()

	// checking the form data
	if data.ValidateTransaction(&form.Validator, *transaction); !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "transaction-create.tmpl")
		return
	}

	err = app.models.TransactionModel.Insert(transaction)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.background(func() {

		// TODO -> update the mailData for the command-receipt mail
		mailData := map[string]any{
			"clientID":   transaction.Client.ID,
			"first-name": transaction.Client.FirstName,
			"last-name":  transaction.Client.LastName,
		}

		err = app.mailer.Send(transaction.Client.Email, "command-receipt.tmpl", mailData)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	// adding notification message
	app.sessionManager.Put(r.Context(), "flash", "Transaction has been created successfully")

	http.Redirect(w, r, fmt.Sprintf("/car-catalog/%d", transaction.ID), http.StatusSeeOther)

}

func (app *application) updateTransaction(w http.ResponseWriter, r *http.Request) {
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Phoceene Auto - Transaction"

	// retrieving ID
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// fetching car data
	transaction, err := app.models.TransactionModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// filling the form with car values
	tmplData.Form = newTransactionUpdateForm(transaction)

	// rendering the template
	app.render(w, r, http.StatusOK, "transaction-update.tmpl", tmplData)
}

func (app *application) updateTransactionPost(w http.ResponseWriter, r *http.Request) {
	// retrieving the form data
	form := newTransactionUpdateForm(nil)
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
	transaction, err := app.models.TransactionModel.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking the data from the car
	isEmpty := form.toTransaction(transaction)
	if isEmpty {
		form.AddNonFieldError("at least one field is required")
	}

	// checking the form data
	if data.ValidateTransaction(&form.Validator, *transaction); !form.Valid() {
		app.failedValidationError(w, r, form, &form.Validator, "transaction-update.tmpl")
		return
	}

	// updating the car
	err = app.models.TransactionModel.Update(transaction, true)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.clientError(w, r, http.StatusNotFound)

		default:
			app.serverError(w, r, err)
		}
		return
	}

	// adding notification message
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Transaction %d has been updated successfully!", transaction.ID))

	http.Redirect(w, r, fmt.Sprintf("/transaction/%d", transaction.ID), http.StatusSeeOther)
}

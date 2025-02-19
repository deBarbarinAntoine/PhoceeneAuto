package main

import (
	"net/http"

	"PhoceeneAuto/internal/data"
)

func (app *application) transactionGet(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) createTransaction(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) createTransactionPost(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement

	var err error
	var transaction data.Transaction

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

}

func (app *application) updateTransaction(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

func (app *application) updateTransactionPost(w http.ResponseWriter, r *http.Request) {
	// TODO -> to implement
}

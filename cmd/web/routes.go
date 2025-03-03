package main

import (
	"io/fs"
	"net/http"

	"PhoceeneAuto/ui"

	"github.com/alexedwards/flow"
)

// routes sets up the HTTP routes for the application.
//
// Parameters:
//
//	app - The application instance
//
// Returns:
//
//	http.Handler - A new HTTP handler with all routes set up
func (app *application) routes() http.Handler {

	// setting the files to put in the static handler
	staticFs, err := fs.Sub(ui.StaticFiles, "assets")
	if err != nil {
		panic(err)
	}

	router := flow.New()

	router.NotFound = http.HandlerFunc(app.notFound)                 // error 404 page
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed) // error 405 page

	router.Handle("/static/...", http.StripPrefix("/static/", http.FileServerFS(staticFs)), http.MethodGet) // static files

	router.Use(app.recoverPanic, app.logRequest, commonHeaders, app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	/* #############################################################################
	/*	RESTRICTED
	/* #############################################################################*/

	router.Group(func(group *flow.Mux) {

		group.Use(app.requireAuthentication)

		// USER
		group.HandleFunc("/dashboard", app.dashboard, http.MethodGet) // dashboard page
		group.HandleFunc("/logout", app.logoutPost, http.MethodPost)  // logout route

		router.HandleFunc("/search", app.searchClassic, http.MethodGet) // search page

		// CLIENT
		router.HandleFunc("/client/:id", app.clientGet, http.MethodGet)       // get client
		router.HandleFunc("/client/:id", app.deleteClient, http.MethodDelete) // delete client treatment route

		router.HandleFunc("/client", app.createClient, http.MethodGet)      // create client page
		router.HandleFunc("/client", app.createClientPost, http.MethodPost) // create client treatment route

		router.HandleFunc("/client/:id/update", app.updateClient, http.MethodGet)      // update client page
		router.HandleFunc("/client/:id/update", app.updateClientPost, http.MethodPost) // update client treatment route

		// CAR PRODUCT
		router.HandleFunc("/car-product/:id", app.carProductGet, http.MethodGet)       // get car-product
		router.HandleFunc("/car-product/:id", app.deleteCarProduct, http.MethodDelete) // delete car-product treatment route

		router.HandleFunc("/car-product", app.createCarProduct, http.MethodGet)      // create car-product page
		router.HandleFunc("/car-product", app.createCarProductPost, http.MethodPost) // create car-product treatment route

		router.HandleFunc("/car-product/:id/update", app.updateCarProduct, http.MethodGet)      // update car-product page
		router.HandleFunc("/car-product/:id/update", app.updateCarProductPost, http.MethodPost) // update car-product treatment route

		// CAR CATALOG
		router.HandleFunc("/car-catalog/:id", app.carCatalogGet, http.MethodGet)       // get car-catalog
		router.HandleFunc("/car-catalog/:id", app.deleteCarCatalog, http.MethodDelete) // delete car-catalog treatment route

		router.HandleFunc("/car-catalog", app.createCarCatalog, http.MethodGet)      // create car-catalog page
		router.HandleFunc("/car-catalog", app.createCarCatalogPost, http.MethodPost) // create car-catalog treatment route

		router.HandleFunc("/car-catalog/:id/update", app.updateCarCatalog, http.MethodGet)      // update car page
		router.HandleFunc("/car-catalog/:id/update", app.updateCarCatalogPost, http.MethodPost) // update car treatment route

		// TRANSACTION
		router.HandleFunc("/transaction/:id", app.transactionGet, http.MethodGet)       // get transaction
		router.HandleFunc("/transaction/:id", app.deleteTransaction, http.MethodDelete) // delete transaction treatment route

		router.HandleFunc("/transaction", app.createTransaction, http.MethodGet)      // create transaction page
		router.HandleFunc("/transaction", app.createTransactionPost, http.MethodPost) // create transaction treatment route

		router.HandleFunc("/transaction/:id/update", app.updateTransaction, http.MethodGet)      // update transaction page
		router.HandleFunc("/transaction/:id/update", app.updateTransactionPost, http.MethodPost) // update transaction treatment route

		/* ########################################################################
		/*	ADMIN AREA
		/* ######################################################################## */

		group.Use(app.requireAdmin)

		group.HandleFunc("/user/:id", app.userGet, http.MethodGet)       // get user
		group.HandleFunc("/user/:id", app.deleteUser, http.MethodDelete) // delete user treatment route

		group.HandleFunc("/user", app.createUser, http.MethodGet)      // create user page
		group.HandleFunc("/user", app.createUserPost, http.MethodPost) // create user treatment route

		group.HandleFunc("/user/:id/update", app.updateUser, http.MethodGet)      // update user page
		group.HandleFunc("/user/:id/update", app.updateUserPost, http.MethodPost) // update user treatment route

		group.HandleFunc("/reports", app.reports, http.MethodGet) // reports and statistics page

	})

	/* #############################################################################
	/*	USER ACCESS
	/* #############################################################################*/

	router.HandleFunc("/", app.login, http.MethodGet)           // login page
	router.HandleFunc("/login", app.login, http.MethodGet)      // login page
	router.HandleFunc("/login", app.loginPost, http.MethodPost) // login treatment route

	return router
}

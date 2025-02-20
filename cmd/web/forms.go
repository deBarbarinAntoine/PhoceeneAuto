package main

import (
	"PhoceeneAuto/internal/data"
	"PhoceeneAuto/internal/validator"
)

// newUserCreateForm creates a new user creation form.
//
// Returns:
//
//	*userCreateForm - The created form
func newUserCreateForm() *userCreateForm {
	return &userCreateForm{
		Validator: *validator.New(),
	}
}

// newUserLoginForm creates a new user login form.
//
// Returns:
//
//	*userLoginForm - The created form
func newUserLoginForm() *userLoginForm {
	return &userLoginForm{
		Validator: *validator.New(),
	}
}

// newUserUpdateForm creates a new user update form.
//
// Parameters:
//
//	user - The user data to pre-fill the form with, or nil if no data is available
//
// Returns:
//
//	*userUpdateForm - The created form
func newUserUpdateForm(user *data.User) *userUpdateForm {

	// creating the form
	var form = new(userUpdateForm)

	// filling the form with the data if any
	if user != nil {
		form.ID = &user.ID
		form.Username = &user.Name
		form.Email = &user.Email
		form.Phone = &user.Phone
		form.Street = &user.Address.Street
		form.Complement = &user.Address.Complement
		form.City = &user.Address.City
		form.ZIP = &user.Address.ZIP
		form.Country = &user.Address.Country
		form.Shop = &user.Shop
		form.Status = &user.Status
		form.Role = &user.Role
	}

	// setting the validator
	form.Validator = *validator.New()

	return form
}

// newClientCreateForm creates a new client creation form.
//
// Returns:
//
//	*clientCreateForm - The created form
func newClientCreateForm() *clientCreateForm {
	return &clientCreateForm{
		Validator: *validator.New(),
	}
}

// newClientUpdateForm creates a new client update form.
//
// Parameters:
//
//	client - The client data to pre-fill the form with, or nil if no data is available
//
// Returns:
//
//	*clientUpdateForm - The created form
func newClientUpdateForm(client *data.Client) *clientUpdateForm {

	// creating the form
	var form = new(clientUpdateForm)

	// filling the form with the data if any
	if client != nil {
		form.ID = &client.ID
		form.FirstName = &client.FirstName
		form.LastName = &client.LastName
		form.Email = &client.Email
		form.Phone = &client.Phone
		form.Street = &client.Address.Street
		form.Complement = &client.Address.Complement
		form.City = &client.Address.City
		form.ZIP = &client.Address.ZIP
		form.Country = &client.Address.Country
		form.Status = &client.Status
		form.Shop = &client.Shop
	}

	// setting the validator
	form.Validator = *validator.New()

	return form
}

// newCarCatalogCreateForm creates a new car catalog creation form.
//
// Returns:
//
//	*carCatalogCreateForm - The created form
func newCarCatalogCreateForm() *carCatalogCreateForm {
	return &carCatalogCreateForm{
		Validator: *validator.New(),
	}
}

// newCarCatalogUpdateForm creates a new car catalog update form.
//
// Parameters:
//
//	car - The car catalog data to pre-fill the form with, or nil if no data is available
//
// Returns:
//
//	*carCatalogUpdateForm - The created form
func newCarCatalogUpdateForm(car *data.CarCatalog) *carCatalogUpdateForm {

	// creating the form
	var form = new(carCatalogUpdateForm)

	// filling the form with the data if any
	if car != nil {
		form.ID = &car.CatID
		form.Make = &car.Make
		form.Model = &car.Model
		form.Year = &car.Year
		form.Transmission = &car.Transmission
		form.Fuel1 = &car.Fuel1
		form.Fuel2 = &car.Fuel2
		form.Cylinders = &car.Cylinders
		form.Drive = &car.Drive
		form.ElectricMotor = &car.ElectricMotor
		form.EngineDescriptor = &car.EngineDescriptor
		form.LuggageVolume = &car.LuggageVolume
		form.PassengerVolume = &car.PassengerVolume
		form.SizeClass = &car.SizeClass
		form.BaseModel = &car.BaseModel
	}

	// setting the validator
	form.Validator = *validator.New()

	return form
}

// newCarProductCreateForm creates a new car product creation form.
//
// Returns:
//
//	*carProductCreateForm - The created form
func newCarProductCreateForm() *carProductCreateForm {
	return &carProductCreateForm{
		Validator: *validator.New(),
	}
}

// newCarProductUpdateForm creates a new car product update form.
//
// Parameters:
//
//	car - The car product data to pre-fill the form with, or nil if no data is available
//
// Returns:
//
//	*carProductUpdateForm - The created form
func newCarProductUpdateForm(car *data.CarProduct) *carProductUpdateForm {

	// creating the form
	var form = new(carProductUpdateForm)

	// filling the form with the data if any
	if car != nil {
		form.ID = &car.ID
		form.CatID = &car.CatID
		form.Shop = &car.Shop
		form.Status = &car.Status
		form.OwnerNb = &car.OwnerNb
		form.Color = &car.Color
		form.Kilometers = &car.Kilometers
	}

	// setting the validator
	form.Validator = *validator.New()

	return form
}

// newTransactionCreateForm creates a new transaction creation form.
//
// Returns:
//
//	*transactionCreateForm - The created form
func newTransactionCreateForm() *transactionCreateForm {
	return &transactionCreateForm{
		Validator: *validator.New(),
	}
}

// newTransactionUpdateForm creates a new transaction update form.
//
// Parameters:
//
//	transaction - The transaction data to pre-fill the form with, or nil if no data is available
//
// Returns:
//
//	*transactionUpdateForm - The created form
func newTransactionUpdateForm(transaction *data.Transaction) *transactionUpdateForm {

	// creating the form
	var form = new(transactionUpdateForm)

	// filling the form with the data if any
	if transaction != nil {
		for _, car := range transaction.Cars {
			form.CarsID = append(form.CarsID, car.ID)
		}
		form.ClientID = &transaction.Client.ID
		form.UserID = &transaction.User.ID
		form.Status = &transaction.Status
		form.Leases = transaction.Leases
		form.TotalPrice = &transaction.TotalPrice
	}

	// setting the validator
	form.Validator = *validator.New()

	return form
}

// newSearchForm creates a new search form.
//
// Returns:
//
//	*searchForm - The created form
func newSearchForm() *searchForm {
	return &searchForm{
		Validator: *validator.New(),
	}
}

func (form userUpdateForm) toUser(user *data.User) bool {
	isEmpty := true

	if form.Username != nil {
		isEmpty = false
		user.Name = *form.Username
	}
	if form.Email != nil {
		isEmpty = false
		user.Email = *form.Email
	}
	if form.Phone != nil {
		isEmpty = false
		user.Phone = *form.Phone
	}
	if form.Street != nil {
		isEmpty = false
		user.Address.Street = *form.Street
	}
	if form.Complement != nil {
		isEmpty = false
		user.Address.Complement = *form.Complement
	}
	if form.City != nil {
		isEmpty = false
		user.Address.City = *form.City
	}
	if form.ZIP != nil {
		isEmpty = false
		user.Address.ZIP = *form.ZIP
	}
	if form.Country != nil {
		isEmpty = false
		user.Address.Country = *form.Country
	}
	if form.Status != nil {
		isEmpty = false
		user.Status = *form.Status
	}
	if form.Shop != nil {
		isEmpty = false
		user.Shop = *form.Shop
	}
	if form.Role != nil {
		isEmpty = false
		user.Role = *form.Role
	}

	return isEmpty
}

func (form userCreateForm) toUser() *data.User {
	user := data.EmptyUser()

	if form.Username != nil {
		user.Name = *form.Username
	}
	if form.Email != nil {
		user.Email = *form.Email
	}
	if form.Phone != nil {
		user.Phone = *form.Phone
	}
	if form.Street != nil {
		user.Address.Street = *form.Street
	}
	if form.Complement != nil {
		user.Address.Complement = *form.Complement
	}
	if form.City != nil {
		user.Address.City = *form.City
	}
	if form.ZIP != nil {
		user.Address.ZIP = *form.ZIP
	}
	if form.Country != nil {
		user.Address.Country = *form.Country
	}
	if form.Status != nil {
		user.Status = *form.Status
	}
	if form.Shop != nil {
		user.Shop = *form.Shop
	}
	if form.Role != nil {
		user.Role = *form.Role
	}

	return user
}

func (form clientCreateForm) toClient() *data.Client {
	client := data.EmptyClient()

	if form.FirstName != nil {
		client.FirstName = *form.FirstName
	}
	if form.LastName != nil {
		client.LastName = *form.LastName
	}
	if form.Email != nil {
		client.Email = *form.Email
	}
	if form.Phone != nil {
		client.Phone = *form.Phone
	}
	if form.Street != nil {
		client.Address.Street = *form.Street
	}
	if form.Complement != nil {
		client.Address.Complement = *form.Complement
	}
	if form.City != nil {
		client.Address.City = *form.City
	}
	if form.ZIP != nil {
		client.Address.ZIP = *form.ZIP
	}
	if form.Country != nil {
		client.Address.Country = *form.Country
	}
	if form.Status != nil {
		client.Status = *form.Status
	}
	if form.Shop != nil {
		client.Shop = *form.Shop
	}

	return client
}

func (form clientUpdateForm) toClient(client *data.Client) {

	if form.FirstName != nil {
		client.FirstName = *form.FirstName
	}
	if form.LastName != nil {
		client.LastName = *form.LastName
	}
	if form.Email != nil {
		client.Email = *form.Email
	}
	if form.Phone != nil {
		client.Phone = *form.Phone
	}
	if form.Street != nil {
		client.Address.Street = *form.Street
	}
	if form.Complement != nil {
		client.Address.Complement = *form.Complement
	}
	if form.City != nil {
		client.Address.City = *form.City
	}
	if form.ZIP != nil {
		client.Address.ZIP = *form.ZIP
	}
	if form.Country != nil {
		client.Address.Country = *form.Country
	}
	if form.Status != nil {
		client.Status = *form.Status
	}
	if form.Shop != nil {
		client.Shop = *form.Shop
	}
}

func (form carCatalogCreateForm) toCarCatalog() *data.CarCatalog {
	car := data.EmptyCarCatalog()

	if form.Make != nil {
		car.Make = *form.Make
	}
	if form.Model != nil {
		car.Model = *form.Model
	}
	if form.Cylinders != nil {
		car.Cylinders = *form.Cylinders
	}
	if form.Drive != nil {
		car.Drive = *form.Drive
	}
	if form.EngineDescriptor != nil {
		car.EngineDescriptor = *form.EngineDescriptor
	}
	if form.Fuel1 != nil {
		car.Fuel1 = *form.Fuel1
	}
	if form.Fuel2 != nil {
		car.Fuel2 = *form.Fuel2
	}
	if form.LuggageVolume != nil {
		car.LuggageVolume = *form.LuggageVolume
	}
	if form.PassengerVolume != nil {
		car.PassengerVolume = *form.PassengerVolume
	}
	if form.Transmission != nil {
		car.Transmission = *form.Transmission
	}
	if form.SizeClass != nil {
		car.SizeClass = *form.SizeClass
	}
	if form.Year != nil {
		car.Year = *form.Year
	}
	if form.ElectricMotor != nil {
		car.ElectricMotor = *form.ElectricMotor
	}
	if form.BaseModel != nil {
		car.BaseModel = *form.BaseModel
	}

	return car
}

func (form carCatalogUpdateForm) toCarCatalog(car *data.CarCatalog) bool {
	isEmpty := true

	if form.Make != nil {
		isEmpty = false
		car.Make = *form.Make
	}
	if form.Model != nil {
		isEmpty = false
		car.Model = *form.Model
	}
	if form.Cylinders != nil {
		isEmpty = false
		car.Cylinders = *form.Cylinders
	}
	if form.Drive != nil {
		isEmpty = false
		car.Drive = *form.Drive
	}
	if form.EngineDescriptor != nil {
		isEmpty = false
		car.EngineDescriptor = *form.EngineDescriptor
	}
	if form.Fuel1 != nil {
		isEmpty = false
		car.Fuel1 = *form.Fuel1
	}
	if form.Fuel2 != nil {
		isEmpty = false
		car.Fuel2 = *form.Fuel2
	}
	if form.LuggageVolume != nil {
		isEmpty = false
		car.LuggageVolume = *form.LuggageVolume
	}
	if form.PassengerVolume != nil {
		isEmpty = false
		car.PassengerVolume = *form.PassengerVolume
	}
	if form.Transmission != nil {
		isEmpty = false
		car.Transmission = *form.Transmission
	}
	if form.SizeClass != nil {
		isEmpty = false
		car.SizeClass = *form.SizeClass
	}
	if form.Year != nil {
		isEmpty = false
		car.Year = *form.Year
	}
	if form.ElectricMotor != nil {
		isEmpty = false
		car.ElectricMotor = *form.ElectricMotor
	}
	if form.BaseModel != nil {
		isEmpty = false
		car.BaseModel = *form.BaseModel
	}
	return isEmpty
}

func (form carProductCreateForm) toCarProduct() *data.CarProduct {
	car := data.EmptyCarProduct()

	if form.OwnerNb != nil {
		car.OwnerNb = *form.OwnerNb
	}
	if form.Color != nil {
		car.Color = *form.Color
	}
	if form.Price != nil {
		car.Price = *form.Price
	}
	if form.Shop != nil {
		car.Shop = *form.Shop
	}
	if form.CatID != nil {
		car.CatID = *form.CatID
	}

	return car
}

func (form carProductUpdateForm) toCarProduct(car *data.CarProduct) bool {
	isEmpty := true

	if form.OwnerNb != nil {
		isEmpty = false
		car.OwnerNb = *form.OwnerNb
	}
	if form.Color != nil {
		isEmpty = false
		car.Color = *form.Color
	}
	if form.Price != nil {
		isEmpty = false
		car.Price = *form.Price
	}
	if form.Shop != nil {
		isEmpty = false
		car.Shop = *form.Shop
	}
	if form.CatID != nil {
		isEmpty = false
		car.CatID = *form.CatID
	}
	return isEmpty
}

func (form transactionCreateForm) toTransaction() *data.Transaction {
	transaction := data.EmptyTransaction()

	if form.CarsID != nil {
		for _, id := range form.CarsID {
			transaction.Cars = append(transaction.Cars, data.CarProduct{ID: id})
		}
	}
	if form.ClientID != nil {
		transaction.Client.ID = *form.ClientID
	}
	if form.UserID != nil {
		transaction.User.ID = *form.UserID
	}
	if form.Status != nil {
		transaction.Status = *form.Status
	}
	if form.Leases != nil {
		transaction.Leases = form.Leases
	}
	if form.TotalPrice != nil {
		transaction.TotalPrice = *form.TotalPrice
	}

	return transaction
}

func (form transactionUpdateForm) toTransaction(transaction *data.Transaction) {

	if form.CarsID != nil {
		for i, id := range form.CarsID {
			transaction.Cars[i].ID = id
		}
	}
	if form.ClientID != nil {
		transaction.Client.ID = *form.ClientID
	}
	if form.UserID != nil {
		transaction.User.ID = *form.UserID
	}
	if form.Status != nil {
		transaction.Status = *form.Status
	}
	if form.Leases != nil {
		transaction.Leases = form.Leases
	}
	if form.TotalPrice != nil {
		transaction.TotalPrice = *form.TotalPrice
	}
}

func (form searchForm) toSearch(search *data.Search) {
	filled := false

	if form.Search != nil {
		search.Search = form.Search
		filled = true
	}
	if form.Make != nil {
		search.Make = form.Make
		filled = true
	}
	if form.Model != nil {
		search.Model = form.Model
		filled = true
	}
	if form.Year != nil {
		search.Year = form.Year
		filled = true
	}
	if form.PriceMin != nil {
		search.PriceMin = form.PriceMin
		filled = true
	}
	if form.PriceMax != nil {
		search.PriceMax = form.PriceMax
		filled = true
	}
	if form.KmMin != nil {
		search.KmMin = form.KmMin
		filled = true
	}
	if form.KmMax != nil {
		search.KmMax = form.KmMax
		filled = true
	}
	if form.Color != nil {
		search.Color = form.Color
		filled = true
	}
	if form.Transmission != nil {
		search.Transmission = form.Transmission
		filled = true
	}
	if form.Fuel1 != nil {
		search.Fuel1 = form.Fuel1
		filled = true
	}
	if form.Fuel2 != nil {
		search.Fuel2 = form.Fuel2
		filled = true
	}
	if form.SizeClass != nil {
		search.SizeClass = form.SizeClass
		filled = true
	}
	if form.OwnerCount != nil {
		search.OwnerCount = form.OwnerCount
		filled = true
	}
	if form.Shop != nil {
		search.Shop = form.Shop
		filled = true
	}
	if form.Status != nil {
		search.Status = form.Status
		filled = true
	}
	if form.ClientName != nil {
		search.ClientName = form.ClientName
		filled = true
	}
	if form.Email != nil {
		search.Email = form.Email
		filled = true
	}
	if form.Phone != nil {
		search.Phone = form.Phone
		filled = true
	}
	if form.ClientStatus != nil {
		search.ClientStatus = form.ClientStatus
		filled = true
	}
	if form.TransactionID != nil {
		search.TransactionID = form.TransactionID
		filled = true
	}
	if form.UserID != nil {
		search.UserID = form.UserID
		filled = true
	}
	if form.TransactionStatus != nil {
		search.TransactionStatus = form.TransactionStatus
		filled = true
	}
	if form.DateStart != nil {
		search.DateStart = form.DateStart
		filled = true
	}
	if form.DateEnd != nil {
		search.DateEnd = form.DateEnd
		filled = true
	}
	if form.LeaseAmountMin != nil {
		search.LeaseAmountMin = form.LeaseAmountMin
		filled = true
	}
	if form.LeaseAmountMax != nil {
		search.LeaseAmountMax = form.LeaseAmountMax
		filled = true
	}

	if !filled {
		// panic("At least one search parameter must be provided")
	}
}

package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type SearchModel struct {
	db *sql.DB
}

type SearchResult struct {
	Transactions []*Transaction
	Clients      []*Client
	CarCatalogs  []*CarCatalog
	CarProducts  []*CarProduct
}

type Search struct {
	Search       *string  `form:"search"` // Classic search bar input
	Make         *string  `form:"make"`   // Car make
	Model        *string  `form:"model"`  // Car model
	Year         *int     `form:"year,omitempty"`
	PriceMin     *float64 `form:"price_min"`    // Minimum price
	PriceMax     *float64 `form:"price_max"`    // Maximum price
	KmMin        *float64 `form:"km_min"`       // Minimum kilometers driven
	KmMax        *float64 `form:"km_max"`       // Maximum kilometers driven
	Color        *string  `form:"color"`        // Car color
	Transmission *string  `form:"transmission"` // Manual or Automatic
	Fuel1        *string  `form:"fuel1"`        // Gas, Diesel, Electric, Hybrid
	Fuel2        *string  `form:"fuel2"`
	SizeClass    *string  `form:"size_class,omitempty"`
	OwnerCount   *int     `form:"owner_count"` // Number of previous owners
	Shop         *string  `form:"shop"`        // Shop name
	Status       *string  `form:"status"`      // Available, Sold, etc.

	// Client-related fields
	ClientName   *string `form:"client_name"`   // First name / Last name search
	Email        *string `form:"email"`         // Email search
	Phone        *string `form:"phone"`         // Phone number search
	ClientStatus *string `form:"client_status"` // Active, Inactive

	// Transaction-related fields
	TransactionID     *int     `form:"transaction_id"`     // Exact match for transaction ID
	UserID            *int     `form:"user_id"`            // Salesperson ID
	TransactionStatus *string  `form:"transaction_status"` // Pending, Completed, etc.
	DateStart         *string  `form:"date_start"`         // Transaction start date
	DateEnd           *string  `form:"date_end"`           // Transaction end date
	LeaseAmountMin    *float64 `form:"lease_min"`          // Minimum lease amount
	LeaseAmountMax    *float64 `form:"lease_max"`          // Maximum lease amount
}

func (m SearchModel) SearchAll(form Search) (SearchResult, error) {
	var result SearchResult
	query := form.Search

	if query == nil || *query == "" {
		return result, fmt.Errorf("search query cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	searchTerm := "%%" + *query + "%%"

	// Search clients
	clients, err := m.searchClients(ctx, searchTerm)
	if err != nil {
		return result, err
	}
	result.Clients = clients

	// Search cars
	cars, err := m.searchCars(ctx, searchTerm)
	if err != nil {
		return result, err
	}
	result.CarCatalogs = cars

	// Search car products
	carProducts, err := m.searchCarProducts(ctx, searchTerm)
	if err != nil {
		return result, err
	}
	result.CarProducts = carProducts

	// Search transactions
	transactions, err := m.searchTransactions(ctx, searchTerm)
	if err != nil {
		return result, err
	}
	result.Transactions = transactions

	return result, nil
}

func (m SearchModel) searchClients(ctx context.Context, searchTerm string) ([]*Client, error) {
	query := `SELECT id, created_at, updated_at,first_name, last_name,email, phone,status, shop, street, complement,city, zip_code,state,version FROM clients WHERE first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1`
	rows, err := m.db.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search clients: %w", err)
	}
	defer rows.Close()

	var clients []*Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt, &client.FirstName, &client.LastName, &client.Email, &client.Phone, &client.Status, &client.Shop, &client.Address.Street, &client.Address.Complement, &client.Address.City, &client.Address.ZIP, &client.Address.Country, &client.Version); err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}
	return clients, nil
}

func (m SearchModel) searchCars(ctx context.Context, searchTerm string) ([]*CarCatalog, error) {
	query := `SELECT id, make, model, transmission FROM cars_catalog WHERE make ILIKE $1 OR model ILIKE $1 OR transmission ILIKE $1`
	rows, err := m.db.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search cars: %w", err)
	}
	defer rows.Close()

	var cars []*CarCatalog
	for rows.Next() {
		var car CarCatalog
		if err := rows.Scan(&car.CatID, &car.Make, &car.Model, &car.Transmission); err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}
	return cars, nil
}

func (m SearchModel) searchCarProducts(ctx context.Context, searchTerm string) ([]*CarProduct, error) {
	query := `SELECT id,created_at, updated_at, status, kilometers, owner_nb, color, price, shop, version, cat_id FROM car_products WHERE status ILIKE $1 OR color ILIKE $1`
	rows, err := m.db.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search car products: %w", err)
	}
	defer rows.Close()

	var carProducts []*CarProduct
	for rows.Next() {
		var carProduct CarProduct
		if err := rows.Scan(&carProduct.ID, &carProduct.Status, &carProduct.Color, &carProduct.CreatedAt, &carProduct.UpdatedAt, &carProduct.Kilometers, &carProduct.OwnerNb, &carProduct.Price, &carProduct.Shop, &carProduct.Version, &carProduct.CatID); err != nil {
			return nil, err
		}
		carProducts = append(carProducts, &carProduct)
	}
	return carProducts, nil
}

func (m SearchModel) searchTransactions(ctx context.Context, searchTerm string) ([]*Transaction, error) {
	query := `SELECT id, created_at, updated_at, status, lease_amount, version FROM transactions WHERE status ILIKE $1`
	rows, err := m.db.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.Leases, &transaction.Version); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (m SearchModel) AdvancedSearch(form Search) (SearchResult, error) {
	var result SearchResult
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Search clients
	clients, err := m.advancedSearchClients(ctx, form)
	if err != nil {
		return result, err
	}
	result.Clients = clients

	// Search cars
	cars, err := m.advancedSearchCars(ctx, form)
	if err != nil {
		return result, err
	}
	result.CarCatalogs = cars

	// Search car products
	carProducts, err := m.advancedSearchCarProducts(ctx, form)
	if err != nil {
		return result, err
	}
	result.CarProducts = carProducts

	// Search transactions
	transactions, err := m.advancedSearchTransactions(ctx, form)
	if err != nil {
		return result, err
	}
	result.Transactions = transactions

	return result, nil
}

func buildFilterQuery(baseQuery string, filters []string) string {
	if len(filters) > 0 {
		baseQuery += " WHERE " + strings.Join(filters, " AND ")
	}
	return baseQuery
}

func (m SearchModel) advancedSearchClients(ctx context.Context, form Search) ([]*Client, error) {
	query := "SELECT id, first_name, last_name, email FROM clients"
	var filters []string
	var args []interface{}

	if form.ClientName != nil {
		filters = append(filters, "(first_name ILIKE $1 OR last_name ILIKE $1)")
		args = append(args, "%"+*form.ClientName+"%")
	}
	if form.Email != nil {
		filters = append(filters, "email ILIKE $2")
		args = append(args, "%"+*form.Email+"%")
	}
	if form.Phone != nil {
		filters = append(filters, "phone ILIKE $3")
		args = append(args, "%"+*form.Phone+"%")
	}
	if form.ClientStatus != nil {
		filters = append(filters, "status = $4")
		args = append(args, *form.ClientStatus)
	}

	query = buildFilterQuery(query, filters)
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search clients: %w", err)
	}
	defer rows.Close()

	var clients []*Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.FirstName, &client.LastName, &client.Email); err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}
	return clients, nil
}

func (m SearchModel) advancedSearchCars(ctx context.Context, form Search) ([]*CarCatalog, error) {
	query := "SELECT id, make, model, transmission FROM cars_catalog"
	var filters []string
	var args []interface{}

	if form.Make != nil {
		filters = append(filters, "make ILIKE $1")
		args = append(args, "%"+*form.Make+"%")
	}
	if form.Model != nil {
		filters = append(filters, "model ILIKE $2")
		args = append(args, "%"+*form.Model+"%")
	}
	if form.Year != nil {
		filters = append(filters, "year = $3")
		args = append(args, *form.Year)
	}
	if form.Color != nil {
		filters = append(filters, "color ILIKE $4")
		args = append(args, "%"+*form.Color+"%")
	}
	if form.Transmission != nil {
		filters = append(filters, "transmission = $5")
		args = append(args, *form.Transmission)
	}

	query = buildFilterQuery(query, filters)
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search cars: %w", err)
	}
	defer rows.Close()

	var cars []*CarCatalog
	for rows.Next() {
		var car CarCatalog
		if err := rows.Scan(&car.CatID, &car.Make, &car.Model, &car.Transmission); err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}
	return cars, nil
}

func (m SearchModel) advancedSearchTransactions(ctx context.Context, form Search) ([]*Transaction, error) {
	query := "SELECT id, status FROM transactions"
	var filters []string
	var args []interface{}

	if form.TransactionID != nil {
		filters = append(filters, "id = $1")
		args = append(args, *form.TransactionID)
	}
	if form.TransactionStatus != nil {
		filters = append(filters, "status = $2")
		args = append(args, *form.TransactionStatus)
	}
	if form.DateStart != nil {
		filters = append(filters, "date >= $3")
		args = append(args, *form.DateStart)
	}
	if form.DateEnd != nil {
		filters = append(filters, "date <= $4")
		args = append(args, *form.DateEnd)
	}

	query = buildFilterQuery(query, filters)
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Status); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (m SearchModel) advancedSearchCarProducts(ctx context.Context, form Search) ([]*CarProduct, error) {
	query := "SELECT id, status, color FROM car_products"
	var filters []string
	var args []interface{}

	if form.Status != nil {
		filters = append(filters, "status = $1")
		args = append(args, *form.Status)
	}
	if form.Color != nil {
		filters = append(filters, "color ILIKE $2")
		args = append(args, "%"+*form.Color+"%")
	}

	query = buildFilterQuery(query, filters)
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search car products: %w", err)
	}
	defer rows.Close()

	var carProducts []*CarProduct
	for rows.Next() {
		var carProduct CarProduct
		if err := rows.Scan(&carProduct.ID, &carProduct.Status, &carProduct.Color); err != nil {
			return nil, err
		}
		carProducts = append(carProducts, &carProduct)
	}
	return carProducts, nil
}

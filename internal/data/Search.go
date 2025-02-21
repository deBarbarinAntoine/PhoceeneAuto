package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
)

type SearchModel struct {
	db *sql.DB
}

// SearchResult represents a combined result from multiple tables
type SearchResult struct {
	TransactionID int64     `json:"transaction_id,omitempty"`
	ClientID      int64     `json:"client_id,omitempty"`
	ClientName    string    `json:"client_name,omitempty"`
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	Status        string    `json:"status,omitempty"`
	CarID         int64     `json:"car_id,omitempty"`
	Make          string    `json:"make,omitempty"`
	Model         string    `json:"model,omitempty"`
	ModelYear     int       `json:"model_year,omitempty"`
	Price         float64   `json:"price,omitempty"`
	Color         string    `json:"color,omitempty"`
	Shop          string    `json:"shop,omitempty"`
	Kilometers    float64   `json:"kilometers,omitempty"`
	OwnerNb       int       `json:"owner_nb,omitempty"`
	LeaseAmount   []float64 `json:"lease_amount,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
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

func (m SearchModel) SearchAll(form Search) ([]SearchResult, error) {
	// Base query joining transactions, clients, cars, and car_products
	query := `
		SELECT 
    t.id AS transaction_id, 
    cl.id AS client_id, 
    (cl.first_name || ' ' || cl.last_name) AS client_name, 
    cl.email, cl.phone, t.status, 
    cp.id AS car_id, c.make, c.model, c.model_year, 
    cp.price, cp.color, cp.shop, 
    cp.kilometers, cp.owner_nb, 
    t.lease_amount, t.created_at, t.updated_at
FROM transactions t
JOIN clients cl ON t.client_id = cl.id
JOIN car_products_transactions cpt ON t.id = cpt.transaction_id
JOIN car_products cp ON cpt.car_product_id = cp.id
JOIN cars_catalog c ON cp.cat_id = c.id
WHERE 1=1
`

	// Store conditions and parameters
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Client filters
	if form.ClientName != nil {
		conditions = append(conditions, fmt.Sprintf("(cl.first_name ILIKE $%d OR cl.last_name ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, "%"+*form.ClientName+"%", "%"+*form.ClientName+"%")
		argIndex += 2
	}
	if form.Email != nil {
		conditions = append(conditions, fmt.Sprintf("cl.email ILIKE $%d", argIndex))
		args = append(args, "%"+*form.Email+"%")
		argIndex++
	}
	if form.Phone != nil {
		conditions = append(conditions, fmt.Sprintf("cl.phone ILIKE $%d", argIndex))
		args = append(args, "%"+*form.Phone+"%")
		argIndex++
	}
	if form.ClientStatus != nil {
		conditions = append(conditions, fmt.Sprintf("cl.status ILIKE $%d", argIndex))
		args = append(args, "%"+*form.ClientStatus+"%")
		argIndex++
	}

	// Transaction filters
	if form.TransactionID != nil {
		conditions = append(conditions, fmt.Sprintf("t.id = $%d", argIndex))
		args = append(args, *form.TransactionID)
		argIndex++
	}
	if form.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("t.user_id = $%d", argIndex))
		args = append(args, *form.UserID)
		argIndex++
	}
	if form.TransactionStatus != nil {
		conditions = append(conditions, fmt.Sprintf("t.status ILIKE $%d", argIndex))
		args = append(args, "%"+*form.TransactionStatus+"%")
		argIndex++
	}
	if form.DateStart != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= $%d", argIndex))
		args = append(args, *form.DateStart)
		argIndex++
	}
	if form.DateEnd != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= $%d", argIndex))
		args = append(args, *form.DateEnd)
		argIndex++
	}

	// Car filters
	if form.Make != nil {
		conditions = append(conditions, fmt.Sprintf("c.make ILIKE $%d", argIndex))
		args = append(args, "%"+*form.Make+"%")
		argIndex++
	}
	if form.Model != nil {
		conditions = append(conditions, fmt.Sprintf("c.model ILIKE $%d", argIndex))
		args = append(args, "%"+*form.Model+"%")
		argIndex++
	}
	if form.Year != nil {
		conditions = append(conditions, fmt.Sprintf("c.model_year = $%d", argIndex))
		args = append(args, *form.Year)
		argIndex++
	}
	if form.PriceMin != nil {
		conditions = append(conditions, fmt.Sprintf("cp.price >= $%d", argIndex))
		args = append(args, *form.PriceMin)
		argIndex++
	}
	if form.PriceMax != nil {
		conditions = append(conditions, fmt.Sprintf("cp.price <= $%d", argIndex))
		args = append(args, *form.PriceMax)
		argIndex++
	}
	if form.Color != nil {
		conditions = append(conditions, fmt.Sprintf("cp.color ILIKE $%d", argIndex))
		args = append(args, "%"+*form.Color+"%")
		argIndex++
	}

	// Apply conditions
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Order by latest transactions
	query += " ORDER BY t.created_at DESC"

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	results, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer results.Close()
	var data []SearchResult
	for results.Next() {
		var result SearchResult
		var leaseAmount pq.Float64Array

		err := results.Scan(
			&result.TransactionID, &result.ClientID, &result.ClientName,
			&result.Email, &result.Phone, &result.Status,
			&result.CarID, &result.Make, &result.Model, &result.ModelYear,
			&result.Price, &result.Color, &result.Shop,
			&result.Kilometers, &result.OwnerNb,
			&leaseAmount, &result.CreatedAt, &result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result.LeaseAmount = leaseAmount
		data = append(data, result)
	}

	// Check for errors after iterating
	if err = results.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return data, nil
}

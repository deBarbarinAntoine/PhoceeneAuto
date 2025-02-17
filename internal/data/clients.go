package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	
	"PhoceeneAuto/internal/validator"
)

var (
	ClientStatus = struct {
		ACTIVE   string
		CONFLICT string
		DELETED  string
	}{
		ACTIVE:   "ACTIVE",
		CONFLICT: "CONFLICT",
		DELETED:  "DELETED",
	}
)

type Client struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Status    string
	Address   Address
	Shop      string
	Version   int
}

func ValidateClient(v *validator.Validator, client *Client) {
	
	v.StringCheck(client.FirstName, 2, 30, true, "first-name")
	v.StringCheck(client.LastName, 2, 40, true, "last-name")
	
	ValidateEmail(v, client.Email)
	
	client.Address.validate(v)
}

type ClientModel struct {
	db *sql.DB
}

func (m ClientModel) Insert(client *Client) error {
	
	// creating the query
	query := `
		INSERT INTO clients (first_name, last_name, email, phone, status, shop, street, complement, city, zip_code, state)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, version;`
	
	// setting the arguments
	args := []any{client.FirstName, client.LastName, client.Email, client.Phone, ClientStatus.ACTIVE, client.Shop}
	client.Address.toSQL(args)
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	
	// executing the query
	err = stmt.QueryRowContext(ctx, args...).Scan(&client.ID, &client.CreatedAt, &client.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "clients_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	
	return nil
}

func (m ClientModel) Update(client *Client) error {
	
	// creating the query
	query := `
		UPDATE clients
		SET first_name = $1,
		    last_name = $2,
		    email = $3,
		    phone = $4,
		    status = $5,
		    shop = $6,
		    street = $9,
		    complement = $10,
		    city = $11,
		    zip_code = $12,
		    state = $13
		    updated_at = CURRENT_DATE,
		    version = version + 1,
		WHERE id = $7 AND version = $8
		RETURNING version;`
	
	// setting the arguments
	args := []any{
		client.FirstName,
		client.LastName,
		client.Email,
		client.Phone,
		client.Status,
		client.Shop,
		client.ID,
		client.Version,
	}
	client.Address.toSQL(args)
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	
	// executing the query
	err = stmt.QueryRowContext(ctx, args...).Scan(&client.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "clients_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	
	return nil
}

func (m ClientModel) Delete(client *Client) error {
	
	// creating the query
	query := `
		UPDATE clients
		SET status = $1,
		    deleted_at = CURRENT_DATE,
		    updated_at = CURRENT_DATE,
		    version = version + 1,
		WHERE id = $2 AND version = $3
		RETURNING version;`
	
	// setting the arguments
	args := []any{
		ClientStatus.DELETED,
		client.ID,
		client.Version,
	}
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	
	// executing the query
	err = stmt.QueryRowContext(ctx, args...).Scan(&client.Version)
	if err != nil {
		return err
	}
	
	return nil
}

func (m ClientModel) DeleteExpired() error {
	
	// generating the query
	query := `
		DELETE FROM clients c
		WHERE c.status = $1 AND (c.deleted_at < CURRENT_DATE - '3 years'::interval);`
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	
	// executing the query
	_, err = stmt.ExecContext(ctx, ClientStatus.DELETED)
	if err != nil {
		return fmt.Errorf("failed to delete expired users: %w", err)
	}
	
	return nil
}

func (m ClientModel) GetByID(id uint) (*Client, error) {
	
	// creating the query
	query := `
		SELECT id, created_at, updated_at,
		       first_name, last_name,
		       email, phone,
		       status, shop,
		       street, complement,
		       city, zip_code,
		       state,
		       version
		FROM clients
		WHERE id = $1;`
	
	// setting the client variable
	var client Client
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	
	// executing the query
	err = stmt.QueryRowContext(ctx, id).Scan(
		&client.ID,
		&client.CreatedAt,
		&client.UpdatedAt,
		&client.FirstName,
		&client.LastName,
		&client.Email,
		&client.Phone,
		&client.Status,
		&client.Shop,
		&client.Address.Street,
		&client.Address.Complement,
		&client.Address.City,
		&client.Address.ZIP,
		&client.Address.State,
		&client.Version,
	)
	
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	
	return &client, nil
}

func (m ClientModel) GetByEmail(email string) (*Client, error) {
	
	// creating the query
	query := `
		SELECT id, created_at, updated_at,
		       first_name, last_name,
		       email, phone,
		       status, shop,
		       street, complement,
		       city, zip_code,
		       state,
		       version
		FROM clients
		WHERE email = $1;`
	
	// setting the client variable
	var client Client
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	
	// executing the query
	err = stmt.QueryRowContext(ctx, email).Scan(
		&client.ID,
		&client.CreatedAt,
		&client.UpdatedAt,
		&client.FirstName,
		&client.LastName,
		&client.Email,
		&client.Phone,
		&client.Status,
		&client.Shop,
		&client.Address.Street,
		&client.Address.Complement,
		&client.Address.City,
		&client.Address.ZIP,
		&client.Address.State,
		&client.Version,
	)
	
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	
	return &client, nil
}

func (m ClientModel) Search(search string, filters *Filters) ([]*Client, Metadata, error) {
	
	// creating the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER,
		       id, created_at, updated_at,
		       first_name, last_name,
		       email, phone,
		       status, shop,
		       street, complement,
		       city, zip_code,
		       state,
		       version
		FROM clients
		WHERE (to_tsvector('simple', first_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (to_tsvector('simple', last_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (to_tsvector('simple', email) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (to_tsvector('simple', phone) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3;`, filters.sortColumn(), filters.sortDirection())
	
	// setting the arguments
	args := []any{search, filters.limit(), filters.offset()}
	
	// setting the variables
	totalRecords := 0
	var clients []*Client
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer stmt.Close()
	
	// executing the query
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	
	// scanning for values
	for rows.Next() {
		var client Client
		
		err := rows.Scan(
			&totalRecords,
			&client.ID,
			&client.CreatedAt,
			&client.UpdatedAt,
			&client.FirstName,
			&client.LastName,
			&client.Email,
			&client.Phone,
			&client.Status,
			&client.Shop,
			&client.Address.Street,
			&client.Address.Complement,
			&client.Address.City,
			&client.Address.ZIP,
			&client.Address.State,
			&client.Version,
		)
		
		if err != nil {
			return nil, Metadata{}, err
		}
		
		// adding the client to the list of matching clients
		clients = append(clients, &client)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	
	// getting the metadata
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	
	return clients, metadata, nil
}

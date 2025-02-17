package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	
	"github.com/lib/pq"
)

var (
	TransactionStatus = struct {
		PROCESSING string
		ONGOING    string
		DONE       string
	}{
		PROCESSING: "PROCESSING",
		ONGOING:    "ONGOING",
		DONE:       "DONE",
	}
)

type Transaction struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Cars       []CarProduct // in the join table
	Client     Client       // client_id
	User       User         // user_id
	Status     string
	Leases     []float32
	TotalPrice float32 // not in the database
	Version    int
}

// TODO -> Transaction check fields

type TransactionModel struct {
	db *sql.DB
}

func (m TransactionModel) Insert(transaction *Transaction) error {
	
	// creating the query
	query := `
		INSERT INTO transactions (client_id, user_id, status, leases)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version;`
	
	// setting the arguments
	args := []any{transaction.Client.ID, transaction.User.ID, transaction.Status, pq.Array(transaction.Leases)}
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the transaction
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// preparing the first query
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	
	// executing the first query
	err = stmt.QueryRowContext(ctx, args...).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.Version)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	
	// looping through the cars ordered
	for _, cars := range transaction.Cars {
		
		// creating the join table query
		query = `
		INSERT INTO car_products_transactions (car_product_id, transaction_id)
		VALUES ($1, $2);`
		
		// setting the arguments
		args = []any{cars.ID, transaction.ID}
		
		// preparing the join table query
		stmt, err = tx.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to prepare query: %w", err)
		}
		
		// executing the join table query
		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return fmt.Errorf("failed to insert car_products_transactions: %w", err)
		}
		
		// closing the statement before looping through the next car
		err = stmt.Close()
		if err != nil {
			return err
		}
	}
	
	// committing the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (m TransactionModel) Update(transaction *Transaction, hasChangedProducts bool) error {
	
	// creating the query
	query := `
		UPDATE transactions
		SET client_id = $1, user_id = $2, status = $3, leases = $4,
		    updated_at = CURRENT_DATE, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version;`
	
	// setting the arguments
	args := []any{
		transaction.Client.ID,
		transaction.User.ID,
		transaction.Status,
		pq.Array(transaction.Leases),
		transaction.ID,
		transaction.Version,
	}
	
	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// preparing the transaction
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// preparing the query
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	
	// executing the query
	err = stmt.QueryRowContext(ctx, args...).Scan(&transaction.Version)
	if err != nil {
		return err
	}
	
	// checking if the join table need updates
	if hasChangedProducts {
		
		// dropping the current carProducts in the transaction
		// creating the query
		query = `
			DELETE FROM car_products_transactions WHERE transaction_id = $1;`
		
		// preparing the query
		stmt, err = tx.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to prepare query: %w", err)
		}
		defer stmt.Close()
		
		_, err = stmt.ExecContext(ctx, transaction.ID)
		if err != nil {
			return fmt.Errorf("failed to clear the car_products_transactions: %w", err)
		}
		
		// looping through the cars ordered to add them afresh
		for _, cars := range transaction.Cars {
			
			// creating the join table query
			query = `
				INSERT INTO car_products_transactions (car_product_id, transaction_id)
				VALUES ($1, $2);`
			
			// setting the arguments
			args = []any{cars.ID, transaction.ID}
			
			// preparing the join table query
			stmt, err = tx.PrepareContext(ctx, query)
			if err != nil {
				return fmt.Errorf("failed to prepare query: %w", err)
			}
			
			// executing the join table query
			_, err = stmt.ExecContext(ctx, args...)
			if err != nil {
				return fmt.Errorf("failed to insert car_products_transactions: %w", err)
			}
			
			// closing the statement before looping through the next car
			err = stmt.Close()
			if err != nil {
				return fmt.Errorf("error closing the statement while updating car_products_transactions: %w", err)
			}
		}
	}
	
	// committing the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (m TransactionModel) Delete(transaction *Transaction) error {
	
	// creating the query
	query := `
		DELETE FROM transactions
		WHERE id = $1 AND version = $2;`
	
	// setting the arguments
	args := []any{
		transaction.ID,
		transaction.Version,
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
	_, err = stmt.ExecContext(ctx, args...)
	
	return err
}

func (m TransactionModel) GetByID(id uint) (*Transaction, error) {
	
	// creating the query
	query := `
		SELECT t.id, t.created_at, t.updated_at, t.status, t.leases, t.version,
		       u.id, u.created_at, u.updated_at, u.name, u.email, u.status, u.version,
		       c.id, c.created_at, c.updated_at,
		       c.first_name, c.last_name,
		       c.email, c.phone,
		       c.status, c.shop,
		       c.street, c.complement,
		       c.city, c.zip_code,
		       c.state,
		       c.version,
		       cp.id, cp.created_at, cp.updated_at,
		       cp.status, cp.kilometers, cp.owner_nb, cp.color, cp.price, cp.shop,
		       cp.version,
		       cp.cat_id,
		       cc.id, cc.created_at, cc.updated_at,
		       cc.make, cc.model,
		       cc.cylinders, cc.drive, cc.engine_descriptor,
		       cc.fuel1, cc.fuel2,
		       cc.luggage_volume, cc.passenger_volume,
		       cc.transmission,
		       cc.size_class,
		       cc.model_year,
		       cc.electric_motor,
		       cc.base_model,
		       cc.version
		FROM transactions t
		INNER JOIN users u ON u.id = t.user_id
		INNER JOIN clients c ON c.id = t.client_id
		INNER JOIN car_products_transactions cpt ON cpt.transaction_id = t.id
		INNER JOIN car_products cp ON cp.id = cpt.car_product_id
		INNER JOIN cars_catalog cc ON cp.id = cc.car_product_id
		WHERE id = $1;`
	
	// setting the variables
	var transaction Transaction
	var cars []CarProduct
	
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
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	
	for rows.Next() {
		var car CarProduct
		
		err = rows.Scan(
			// transaction
			&transaction.ID,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.Status,
			pq.Array(&transaction.Leases),
			&transaction.Version,
			
			// user
			&transaction.User.ID,
			&transaction.User.CreatedAt,
			&transaction.User.UpdatedAt,
			&transaction.User.Name,
			&transaction.User.Email,
			&transaction.User.Password.hash,
			&transaction.User.Status,
			&transaction.User.Version,
			
			// client
			&transaction.Client.ID,
			&transaction.Client.CreatedAt,
			&transaction.Client.UpdatedAt,
			&transaction.Client.FirstName,
			&transaction.Client.LastName,
			&transaction.Client.Email,
			&transaction.Client.Phone,
			&transaction.Client.Status,
			&transaction.Client.Shop,
			&transaction.Client.Address.Street,
			&transaction.Client.Address.Complement,
			&transaction.Client.Address.City,
			&transaction.Client.Address.ZIP,
			&transaction.Client.Address.State,
			&transaction.Client.Version,
			
			// car_product
			&car.ID,
			&car.CreatedAt,
			&car.UpdatedAt,
			&car.Status,
			&car.Kilometers,
			&car.OwnerNb,
			&car.Color,
			&car.Price,
			&car.Shop,
			&car.Version,
			&car.CatID,
			&car.CatCreatedAt,
			&car.CatUpdatedAt,
			&car.Make,
			&car.Model,
			&car.Cylinders,
			&car.Drive,
			&car.EngineDescriptor,
			&car.Fuel1,
			&car.Fuel2,
			&car.LuggageVolume,
			&car.PassengerVolume,
			&car.Transmission,
			&car.SizeClass,
			&car.Year,
			&car.ElectricMotor,
			&car.BaseModel,
			&car.CatVersion,
		)
		
		if err != nil {
			return nil, err
		}
		
		cars = append(cars, car)
	}
	
	// calculating the total price
	for _, cars := range transaction.Cars {
		transaction.TotalPrice += cars.Price
	}
	
	return &transaction, nil
}

func (m TransactionModel) Search(search string, filters *Filters) ([]*Transaction, Metadata, error) {
	
	// TODO -> implement search transaction
	
	return nil, Metadata{}, nil
}

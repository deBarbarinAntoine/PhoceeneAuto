package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"PhoceeneAuto/internal/validator"
)

type CarCatalog struct {
	CatID            int       // id
	CatCreatedAt     time.Time // created_at
	CatUpdatedAt     time.Time // updated_at
	Make             string
	Model            string
	Cylinders        int
	Drive            string
	EngineDescriptor string
	Fuel1            string
	Fuel2            string
	LuggageVolume    float32
	PassengerVolume  float32
	Transmission     string
	SizeClass        string
	Year             int // model_year
	ElectricMotor    float32
	BaseModel        string
	CatVersion       int // version
}

func EmptyCarCatalog() *CarCatalog {
	return &CarCatalog{}
}

func ValidateCarCatalog(v *validator.Validator, car CarCatalog) {

	v.Check(car.Year > 1980, "year", "must be greater than 1980")

	v.StringCheck(car.Make, 2, 50, true, "make")

	v.StringCheck(car.Model, 2, 50, true, "model")

	v.StringCheck(car.Transmission, 2, 50, true, "transmission")

	v.StringCheck(car.Drive, 2, 50, true, "drive")

	v.StringCheck(car.EngineDescriptor, 2, 50, false, "engine_descriptor")

	v.StringCheck(car.Fuel2, 2, 50, false, "fuel_2")

	v.StringCheck(car.SizeClass, 2, 50, false, "size_class")

	v.StringCheck(car.BaseModel, 2, 50, false, "base_model")

	v.Check(car.Cylinders >= 0, "cylinders", "must be equal or greater than 0")

	v.Check(car.LuggageVolume >= 0, "luggage_volume", "must be equal or greater than 0")

	v.Check(car.PassengerVolume >= 0, "passenger_volume", "must be equal or greater than 0")

	v.Check(car.ElectricMotor >= 0, "electric_motor", "must be equal or greater than 0")
}

type CarCatalogModel struct {
	db *sql.DB
}

func (m CarCatalogModel) Insert(car *CarCatalog) error {

	// creating the query
	query := `
		INSERT INTO cars_catalog
		    (make, model,
		     cylinders, drive, engine_descriptor,
		     fuel1, fuel2,
		     luggage_volume, passenger_volume,
		     transmission,
		     size_class,
		     model_year,
		     electric_motor,
		     base_model)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, version;`

	// setting the arguments
	args := []any{
		car.Make, car.Model,
		car.Cylinders, car.Drive, car.EngineDescriptor,
		car.Fuel1, car.Fuel2,
		car.LuggageVolume, car.PassengerVolume,
		car.Transmission,
		car.SizeClass,
		car.Year,
		car.ElectricMotor,
		car.BaseModel,
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
	err = stmt.QueryRowContext(ctx, args...).Scan(&car.CatID, &car.CatCreatedAt, &car.CatVersion)
	if err != nil {
		// Check if error message contains the trigger exception message
		if strings.Contains(err.Error(), "Exact duplicate row: This car already exists in the catalog") {
			return ErrDuplicateCarCatalog
		}
		return err
	}
	return nil
}

func (m CarCatalogModel) Update(car *CarCatalog) error {

	// creating the query
	query := `
		UPDATE cars_catalog
		SET make = $1,
		    model = $2,
		    cylinders = $3,
		    drive = $4,
		    engine_descriptor = $5,
		    fuel1 = $6,
		    fuel2 = $7,
		    luggage_volume = $8,
		    passenger_volume = $9,
		    transmission = $10,
		    size_class = $11,
		    model_year = $12,
		    electric_motor = $13,
		    base_model = $14,
		    updated_at = CURRENT_TIMESTAMP,
		    version = version + 1
		WHERE id = $15 AND version = $16
		RETURNING version;`

	// setting the arguments
	args := []any{
		car.Make,
		car.Model,
		car.Cylinders,
		car.Drive,
		car.EngineDescriptor,
		car.Fuel1,
		car.Fuel2,
		car.LuggageVolume,
		car.PassengerVolume,
		car.Transmission,
		car.SizeClass,
		car.Year,
		car.ElectricMotor,
		car.BaseModel,
		car.CatID,
		car.CatVersion,
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
	err = stmt.QueryRowContext(ctx, args...).Scan(&car.CatVersion)
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

func (m CarCatalogModel) Delete(car *CarCatalog) error {

	// creating the query
	query := `
		DELETE FROM cars_catalog
		WHERE id = $1 AND version = $2;`

	// setting the arguments
	args := []any{
		car.CatID,
		car.CatVersion,
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

func (m CarCatalogModel) GetByID(id int) (*CarCatalog, error) {

	// creating the query
	query := `
		SELECT id, created_at, updated_at,
		       make, model,
		       cylinders, drive, engine_descriptor,
		       fuel1, fuel2,
		       luggage_volume, passenger_volume,
		       transmission,
		       size_class,
		       model_year,
		       electric_motor,
		       base_model,
		       version
		FROM cars_catalog
		WHERE id = $1;`

	// setting the car variable
	var car CarCatalog

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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &car, nil
}

func (m CarCatalogModel) Search(search string, filters *Filters) ([]*CarCatalog, Metadata, error) {

	// TODO -> update the method and query to accept specific filters:
	// (make, model, cylinders, drive, fuel, transmission, size_class, model_year, etc.)

	// creating the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),
		       id, created_at, updated_at,
		       make, model,
		       cylinders, drive, engine_descriptor,
		       fuel1, fuel2,
		       luggage_volume, passenger_volume,
		       transmission,
		       size_class,
		       model_year,
		       electric_motor,
		       base_model,
		       version
		FROM cars_catalog
		WHERE (to_tsvector('simple', make) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (to_tsvector('simple', model) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (to_tsvector('simple', base_model) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3;`, filters.sortColumn(), filters.sortDirection())

	// setting the arguments
	args := []any{search, filters.limit(), filters.offset()}

	// setting the variables
	totalRecords := 0
	var cars []*CarCatalog

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
		var car CarCatalog

		err := rows.Scan(
			&totalRecords,
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
			return nil, Metadata{}, err
		}

		// adding the car to the list of matching CarsCatalog
		cars = append(cars, &car)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// getting the metadata
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return cars, metadata, nil
}

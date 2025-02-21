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
	CarStatus = struct {
		COMING     string
		ONSALE     string
		PROCESSING string
		SOLD       string
	}{
		COMING:     "COMING",
		ONSALE:     "ONSALE",
		PROCESSING: "PROCESSING",
		SOLD:       "SOLD",
	}
)

type CarProduct struct {
	ID         int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     string
	Kilometers float32
	OwnerNb    int
	Color      string
	Price      float32
	Shop       string
	Version    int
	CarCatalog // cat_id
}

func EmptyCarProduct() *CarProduct {
	return &CarProduct{
		Status:     CarStatus.ONSALE,
		Shop:       Shop.HEADQUARTERS,
		CarCatalog: CarCatalog{},
	}
}

func ValidateCarProduct(v *validator.Validator, car CarProduct) {
	v.Check(validator.PermittedValue(car.Status, CarStatus.COMING, CarStatus.ONSALE, CarStatus.PROCESSING, CarStatus.SOLD), "status", fmt.Sprintf("invalid status %s", car.Status))

	v.Check(validator.PermittedValue(car.Shop, Shop.HEADQUARTERS), "shop", fmt.Sprintf("invalid shop %s", car.Shop))

	v.Check(car.Price > 0, "price", "must be greater than 0")

	v.CheckID(car.CatID, "cat_id")
}

type CarProductModel struct {
	db *sql.DB
}

func (m CarProductModel) Insert(car *CarProduct) error {

	// creating the query
	query := `
		INSERT INTO car_products
		    (status, kilometers, owner_nb, color, price, shop, cat_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, version;`

	// setting the arguments
	args := []any{
		car.Status,
		car.Kilometers,
		car.OwnerNb,
		car.Color,
		car.Price,
		car.Shop,
		car.CatID,
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
	err = stmt.QueryRowContext(ctx, args...).Scan(&car.ID, &car.CreatedAt, &car.Version)
	if err != nil {
		return err
	}

	return nil
}

func (m CarProductModel) Update(car *CarProduct) error {

	// creating the query
	query := `
		UPDATE car_products
		SET status = $1,
		    kilometers = $2,
		    owner_nb = $3,
		    color = $4,
		    price = $5,
		    shop = $6,
		    cat_id = $7,
		    updated_at = CURRENT_TIMESTAMP,
		    version = version + 1,
		WHERE id = $8 AND version = $9
		RETURNING version;`

	// setting the arguments
	args := []any{
		car.Status,
		car.Kilometers,
		car.OwnerNb,
		car.Color,
		car.Price,
		car.Shop,
		car.CatID,
		car.ID,
		car.Version,
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
		return err
	}

	return nil
}

func (m CarProductModel) Delete(car *CarProduct) error {

	// creating the query
	query := `
		DELETE FROM car_products
		WHERE id = $1 AND version = $2;`

	// setting the arguments
	args := []any{
		car.ID,
		car.Version,
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

// ExistsCatID checks if a CarCatalog is bound by a CarProduct (may do that in the migrations too for security).
func (m CarProductModel) ExistsCatID(catID int) (bool, error) {

	// creating the query
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM car_products
			WHERE cat_id = $1
		);`

	// setting the car variable
	var exists bool

	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return exists, err
	}
	defer stmt.Close()

	// execute the query
	err = stmt.QueryRowContext(ctx, catID).Scan(&exists)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return exists, ErrRecordNotFound
		default:
			return exists, err
		}
	}

	return exists, nil
}

func (m CarProductModel) GetByID(id int) (*CarProduct, error) {

	// creating the query
	query := `
		SELECT cp.id, cp.created_at, cp.updated_at,
		       cp.status, cp.kilometers, cp.owner_nb, cp.color, cp.price, cp.shop,
		       cp.version,
		       cp.cat_id,
		       cc.created_at, cc.updated_at,
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
		FROM car_products cp
		INNER JOIN cars_catalog cc ON cp.cat_id = cc.id
		WHERE cp.id = $1;`

	// setting the car variable
	var car CarProduct
	var carCatalogSql CarCatalogSql

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
		&carCatalogSql.CatID,
		&carCatalogSql.CatCreatedAt,
		&carCatalogSql.CatUpdatedAt,
		&carCatalogSql.Make,
		&carCatalogSql.Model,
		&carCatalogSql.Cylinders,
		&carCatalogSql.Drive,
		&carCatalogSql.EngineDescriptor,
		&carCatalogSql.Fuel1,
		&carCatalogSql.Fuel2,
		&carCatalogSql.LuggageVolume,
		&carCatalogSql.PassengerVolume,
		&carCatalogSql.Transmission,
		&carCatalogSql.SizeClass,
		&carCatalogSql.Year,
		&carCatalogSql.ElectricMotor,
		&carCatalogSql.BaseModel,
		&carCatalogSql.CatVersion,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	car.CarCatalog = carCatalogSql.ToCarCatalog()

	return &car, nil
}

func (m CarProductModel) Search(search string, filters *Filters) ([]*CarProduct, Metadata, error) {

	// TODO -> update the method and query to accept specific filters:
	// (kilometers, status, owner_nb, color, price, shop,
	// make, model, cylinders, drive, fuel, transmission, size_class, model_year, etc.)

	// creating the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),
		       cp.id, cp.created_at, cp.updated_at,
		       cp.status, cp.kilometers, cp.owner_nb, cp.color, cp.price, cp.shop,
		       cp.version,
		       cp.cat_id,
		       cc.created_at, cc.updated_at,
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
		FROM car_products cp
		INNER JOIN cars_catalog cc ON cp.cat_id = cc.id
		WHERE (to_tsvector('simple', cc.make) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (to_tsvector('simple', cc.model) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (to_tsvector('simple', cc.base_model) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3;`, filters.sortColumn(), filters.sortDirection())

	// setting the arguments
	args := []any{search, filters.limit(), filters.offset()}

	// setting the variables
	totalRecords := 0
	var cars []*CarProduct

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
		var car CarProduct
		var carCatalogSql CarCatalogSql

		err := rows.Scan(
			&totalRecords,
			&car.ID,
			&car.CreatedAt,
			&car.UpdatedAt,
			&car.Status,
			&car.Kilometers,
			&car.OwnerNb,
			&car.Color,
			&car.Price,
			&car.Shop,
			&car.ID,
			&car.Version,
			&carCatalogSql.CatID,
			&carCatalogSql.CatCreatedAt,
			&carCatalogSql.CatUpdatedAt,
			&carCatalogSql.Make,
			&carCatalogSql.Model,
			&carCatalogSql.Cylinders,
			&carCatalogSql.Drive,
			&carCatalogSql.EngineDescriptor,
			&carCatalogSql.Fuel1,
			&carCatalogSql.Fuel2,
			&carCatalogSql.LuggageVolume,
			&carCatalogSql.PassengerVolume,
			&carCatalogSql.Transmission,
			&carCatalogSql.SizeClass,
			&carCatalogSql.Year,
			&carCatalogSql.ElectricMotor,
			&carCatalogSql.BaseModel,
			&carCatalogSql.CatVersion,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		car.CarCatalog = carCatalogSql.ToCarCatalog()

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

package main

import (
	"html/template"
	"log/slog"
	"sync"
	"time"

	"PhoceeneAuto/internal/data"
	"PhoceeneAuto/internal/mailer"
	"PhoceeneAuto/internal/validator"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

// config represents the configuration variables for the application.
type config struct {
	port int64
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	smtp struct {
		host     string
		port     int64
		username string
		password string
		sender   string
	}
}

// application represents the application configuration.
type application struct {
	logger         *slog.Logger
	mailer         mailer.Mailer
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	models         data.Models
	config         *config
	wg             *sync.WaitGroup
}

// templateData represents the data structure used in templates.
type templateData struct {
	Title           string
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	Nonce           string
	CSRFToken       string
	ResetToken      string
	Error           struct {
		Title   string
		Message string
	}
	FieldErrors    map[string]string
	NonFieldErrors []string
	User           data.User
	Client         data.Client
	Clients        []*data.Client
	CarCatalog     data.CarCatalog
	CarsCatalog    []*data.CarCatalog
	CarProduct     data.CarProduct
	CarProducts    []*data.CarProduct
	Transaction    data.Transaction
	Transactions   []*data.Transaction
}

// envelope is a data type for JSON responses.
type envelope map[string]any

// userLoginForm represents the form used for user login.
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// userUpdateForm represents the form used for updating a user's information.
type userUpdateForm struct {
	ID                   *int    `form:"id,omitempty"`
	Username             *string `form:"username,omitempty"`
	Email                *string `form:"email,omitempty"`
	Phone                *string `form:"phone,omitempty"`
	Street               *string `form:"street,omitempty"`
	Complement           *string `form:"complement,omitempty"`
	City                 *string `form:"city,omitempty"`
	ZIP                  *string `form:"zip,omitempty"`
	Country              *string `form:"country,omitempty"`
	Password             *string `form:"password,omitempty"`
	NewPassword          *string `form:"new_password,omitempty"`
	ConfirmationPassword *string `form:"confirmation_password,omitempty"`
	Role                 *string `form:"role,omitempty"`
	Status               *string `form:"status,omitempty"`
	Shop                 *string `form:"shop,omitempty"`
	validator.Validator  `form:"-"`
}

// userCreateForm represents the form used for creating a new user.
type userCreateForm struct {
	Username            *string `form:"username"`
	Email               *string `form:"email"`
	Phone               *string `form:"phone,omitempty"`
	Street              *string `form:"street,omitempty"`
	Complement          *string `form:"complement,omitempty"`
	City                *string `form:"city,omitempty"`
	ZIP                 *string `form:"zip,omitempty"`
	Country             *string `form:"country,omitempty"`
	Password            string  `form:"password"`
	ConfirmPassword     string  `form:"confirm_password"`
	Role                *string `form:"role,omitempty"`
	Status              *string `form:"status,omitempty"`
	Shop                *string `form:"shop,omitempty"`
	validator.Validator `form:"-"`
}

// clientCreateForm represents the form used for creating a new client.
type clientCreateForm struct {
	FirstName           *string `form:"first-name,omitempty"`
	LastName            *string `form:"last-name,omitempty"`
	Email               *string `form:"email,omitempty"`
	Phone               *string `form:"phone,omitempty"`
	Street              *string `form:"street,omitempty"`
	Complement          *string `form:"complement,omitempty"`
	City                *string `form:"city,omitempty"`
	ZIP                 *string `form:"zip,omitempty"`
	Country             *string `form:"country,omitempty"`
	Status              *string `form:"status,omitempty"`
	Shop                *string `form:"shop,omitempty"`
	validator.Validator `form:"-"`
}

// clientUpdateForm represents the form used for updating a client's information.
type clientUpdateForm struct {
	ID                  *int    `form:"id,omitempty"`
	FirstName           *string `form:"first-name,omitempty"`
	LastName            *string `form:"last-name,omitempty"`
	Email               *string `form:"email,omitempty"`
	Phone               *string `form:"phone,omitempty"`
	Street              *string `form:"street,omitempty"`
	Complement          *string `form:"complement,omitempty"`
	City                *string `form:"city,omitempty"`
	ZIP                 *string `form:"zip,omitempty"`
	Country             *string `form:"country,omitempty"`
	Status              *string `form:"status,omitempty"`
	Shop                *string `form:"shop,omitempty"`
	validator.Validator `form:"-"`
}

// carCatalogCreateForm represents the form used for creating a new car catalog entry.
type carCatalogCreateForm struct {
	Make                *string  `form:"make,omitempty"`
	Model               *string  `form:"model,omitempty"`
	Cylinders           *int     `form:"cylinders,omitempty"`
	Drive               *string  `form:"drive,omitempty"`
	EngineDescriptor    *string  `form:"engine_descriptor,omitempty"`
	Fuel1               *string  `form:"fuel_1,omitempty"`
	Fuel2               *string  `form:"fuel_2,omitempty"`
	LuggageVolume       *float32 `form:"luggage_volume,omitempty"`
	PassengerVolume     *float32 `form:"passenger_volume,omitempty"`
	Transmission        *string  `form:"transmission,omitempty"`
	SizeClass           *string  `form:"size_class,omitempty"`
	Year                *int     `form:"year,omitempty"`
	ElectricMotor       *float32 `form:"electric_motor,omitempty"`
	BaseModel           *string  `form:"base_model,omitempty"`
	validator.Validator `form:"-"`
}

// carCatalogUpdateForm represents the form used for updating a car catalog entry.
type carCatalogUpdateForm struct {
	ID                  *int     `form:"id,omitempty"`
	Make                *string  `form:"make,omitempty"`
	Model               *string  `form:"model,omitempty"`
	Cylinders           *int     `form:"cylinders,omitempty"`
	Drive               *string  `form:"drive,omitempty"`
	EngineDescriptor    *string  `form:"engine_descriptor,omitempty"`
	Fuel1               *string  `form:"fuel_1,omitempty"`
	Fuel2               *string  `form:"fuel_2,omitempty"`
	LuggageVolume       *float32 `form:"luggage_volume,omitempty"`
	PassengerVolume     *float32 `form:"passenger_volume,omitempty"`
	Transmission        *string  `form:"transmission,omitempty"`
	SizeClass           *string  `form:"size_class,omitempty"`
	Year                *int     `form:"year,omitempty"`
	ElectricMotor       *float32 `form:"electric_motor,omitempty"`
	BaseModel           *string  `form:"base_model,omitempty"`
	validator.Validator `form:"-"`
}

// carProductCreateForm represents the form used for creating a new car product.
type carProductCreateForm struct {
	Status              *string  `form:"status,omitempty"`
	Kilometers          *float32 `form:"kilometers,omitempty"`
	OwnerNb             *int     `form:"owner_nb,omitempty"`
	Color               *string  `form:"color,omitempty"`
	Price               *float32 `form:"price,omitempty"`
	Shop                *string  `form:"shop,omitempty"`
	CatID               *int     `form:"cat_id,omitempty"`
	validator.Validator `form:"-"`
}

// carProductUpdateForm represents the form used for updating a car product.
type carProductUpdateForm struct {
	ID                  *int     `form:"id,omitempty"`
	Status              *string  `form:"status,omitempty"`
	Kilometers          *float32 `form:"kilometers,omitempty"`
	OwnerNb             *int     `form:"owner_nb,omitempty"`
	Color               *string  `form:"color,omitempty"`
	Price               *float32 `form:"price,omitempty"`
	Shop                *string  `form:"shop,omitempty"`
	CatID               *int     `form:"cat_id,omitempty"`
	validator.Validator `form:"-"`
}

// transactionCreateForm represents the form used for creating a new transaction.
type transactionCreateForm struct {
	CarsID              []int     `form:"cars_id,omitempty"`
	ClientID            *int      `form:"client_id,omitempty"`
	UserID              *int      `form:"user_id,omitempty"`
	Status              *string   `form:"status,omitempty"`
	Leases              []float32 `form:"leases,omitempty"`
	TotalPrice          *float32  `form:"total_price,omitempty"`
	validator.Validator `form:"-"`
}

// transactionUpdateForm represents the form used for updating a transaction.
type transactionUpdateForm struct {
	CarsID              []int     `form:"cars_id,omitempty"`
	ClientID            *int      `form:"client_id,omitempty"`
	UserID              *int      `form:"user_id,omitempty"`
	Status              *string   `form:"status,omitempty"`
	Leases              []float32 `form:"leases,omitempty"`
	TotalPrice          *float32  `form:"total_price,omitempty"`
	validator.Validator `form:"-"`
}

// searchForm represents the form used for searching.
type searchForm struct {
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

	validator.Validator `form:"-"` // Keep the validator
}

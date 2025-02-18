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
	Search         string
	CarCatalog     data.CarCatalog
	CarsCatalog    []*data.CarCatalog
	CarProduct     data.CarProduct
	CarProducts    []*data.CarProduct
	Transaction    data.Transaction
	Transactions   []*data.Transaction
}

// envelope data type for JSON responses
type envelope map[string]any

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

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

type userCreateForm struct {
	Username            string  `form:"username"`
	Email               string  `form:"email"`
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

type clientUpdateForm struct {
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

type carCatalogUpdateForm struct {
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

type carProductCreateForm struct {
	Status              *string  `form:"status,omitempty"`
	Kilometers          *float32 `form:"kilometers,omitempty"`
	OwnerNb             *int     `form:"owner_nb,omitempty"`
	Color               *string  `form:"color,omitempty"`
	Price               *float32 `form:"price,omitempty"`
	Shop                *string  `form:"shop,omitempty"`
	validator.Validator `form:"-"`
}

type carProductUpdateForm struct {
	Status              *string  `form:"status,omitempty"`
	Kilometers          *float32 `form:"kilometers,omitempty"`
	OwnerNb             *int     `form:"owner_nb,omitempty"`
	Color               *string  `form:"color,omitempty"`
	Price               *float32 `form:"price,omitempty"`
	Shop                *string  `form:"shop,omitempty"`
	validator.Validator `form:"-"`
}

type transactionCreateForm struct {
	CarsID              []float32 `form:"cars_id,omitempty"`
	ClientID            *int      `form:"client_id,omitempty"`
	UserID              *int      `form:"user_id,omitempty"`
	Status              *string   `form:"status,omitempty"`
	Leases              []float32 `form:"leases,omitempty"`
	TotalPrice          *float32  `form:"total_price,omitempty"`
	validator.Validator `form:"-"`
}

type transactionUpdateForm struct {
	CarsID              []float32 `form:"cars_id,omitempty"`
	ClientID            *int      `form:"client_id,omitempty"`
	UserID              *int      `form:"user_id,omitempty"`
	Status              *string   `form:"status,omitempty"`
	Leases              []float32 `form:"leases,omitempty"`
	TotalPrice          *float32  `form:"total_price,omitempty"`
	validator.Validator `form:"-"`
}

type searchForm struct {
	search *string
	// TODO -> fill the searchForm struct here
	validator.Validator `form:"-"`
}

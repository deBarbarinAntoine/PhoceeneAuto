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
}

// envelope data type for JSON responses
type envelope map[string]any

type contactForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Message             string `form:"message"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userUpdateForm struct {
	Username             *string `form:"username,omitempty"`
	Email                *string `form:"email,omitempty"`
	Password             *string `form:"password,omitempty"`
	NewPassword          *string `form:"new_password,omitempty"`
	ConfirmationPassword *string `form:"confirmation_password,omitempty"`
	Avatar               *string `form:"avatar,omitempty"`
	validator.Validator  `form:"-"`
}

type userCreateForm struct {
	Username            string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	ConfirmPassword     string `form:"confirm_password"`
	validator.Validator `form:"-"`
}

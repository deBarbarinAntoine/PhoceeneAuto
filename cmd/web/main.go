package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"PhoceeneAuto/internal/data"
	"PhoceeneAuto/internal/mailer"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

// main is the entry point of the application.
func main() {

	// setting the configuration variables
	var cfg config

	// generic variables
	flag.Int64Var(&cfg.port, "port", 8080, "HTTP service address")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// PostgreSQL variables
	flag.StringVar(&cfg.db.dsn, "dsn", "", "PostgreSQL Database DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	// SMTP variables
	flag.StringVar(&cfg.smtp.host, "smtp-host", "", "SMTP host")
	flag.Int64Var(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Phoceene Auto <no-reply@phoceene-auto.com", "SMTP sender")

	// cleaning frequency
	frequency := flag.Duration("frequency", time.Hour*2, "expired tokens and unactivated users cleaning frequency")

	flag.Parse()

	// setting the logging level according to the environment
	var opts *slog.HandlerOptions

	if cfg.env == "development" {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	} else {
		opts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// checking the SMTP info
	if cfg.smtp.username == "" || cfg.smtp.password == "" || cfg.smtp.host == "" {
		fmt.Println("SMTP credentials are required")
		os.Exit(1)
	}

	// checking the dsn info
	if cfg.db.dsn == "" {
		logger.Error("dsn is required")
		os.Exit(1)
	}

	// connecting to the database
	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.Error(fmt.Errorf("openDB error: %w", err).Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("DB pool connection opened successfully!")

	// caching the templates
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// initializing the application components
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		mailer:         mailer.New(cfg.smtp.host, int(cfg.smtp.port), cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		sessionManager: sessionManager,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		config:         &cfg,
		models:         data.NewModels(db),
		wg:             new(sync.WaitGroup),
	}

	// Clean deleted clients after a legal period (RGPD) with no timeout
	go app.cleanExpiredDeletedClients(*frequency, time.Hour*0)

	// Running the server
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(1)
}

// openDB opens a new PostgreSQL database connection.
//
// Parameters:
//
//	dsn - The Data Source Name (DSN) used to connect to the database
//
// Returns:
//
//	*sql.DB - A pointer to the opened database connection
//	error - If any error occurs during the process
func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

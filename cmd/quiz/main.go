package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/jsonlog"
	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/model"
	"github.com/margulan-kalykul/JustQuiz/pkg/vcs"
	"github.com/peterbourgon/ff/v3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

// Set version of application corresponding to value of vcs.Version.
var (
	version = vcs.Version()
)

type config struct {
	port       int
	env        string
	fill       bool
	migrations string
	db         struct {
		dsn string
	}
}
type application struct {
	config	config
	models	model.Models
	logger	*jsonlog.Logger
	wg		sync.WaitGroup
}

func main() {
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	var (
		cfg        config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "file://pkg/quiz/migrations", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8081, "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "PostgreSQL DSN")
	)

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	cfg.port = *port
	cfg.env = *env
	cfg.fill = *fill
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	logger.PrintInfo("starting application with configuration", map[string]string{
		"port":       fmt.Sprintf("%d", cfg.port),
		"fill":       fmt.Sprintf("%t", cfg.fill),
		"env":        cfg.env,
		"db":         cfg.db.dsn,
		"migrations": cfg.migrations,
	})

	db, err := openDB(cfg)
	// logger.PrintInfo("", map[string]string{}) // checkpoint
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	// Defer a call to db.Close() so that the connection pool is closed before the main()
	// function exits.
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.AllModels(db),
		logger: logger,
	}
	
	// Call app.server() to start the server.
	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// https://github.com/golang-migrate/migrate?tab=readme-ov-file#use-in-your-go-project
	if cfg.migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.migrations,
			"postgres", driver)
		if err != nil {
			return nil, err
		}
		m.Up()
	}

	return db, nil
}

package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"movieDB/internal/data"
	"movieDB/internal/jsonlog"
	"movieDB/internal/mailer"
	"os"
	"sync"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxOpenConn int
		maxIdleConn int
		maxIdleTime string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgresSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConn, "db-max-open-conn", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConn, "db-max-idle-conn", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "rate limiter max requests per second")
	flag.IntVar(&cfg.limiter.burst, "limited-burst", 4, "rate limiter max burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP Port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "xxx", "SMTP Username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "xxx", "SMTP Password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "William@WillsApp.com", "SMTP sender")

	flag.Parse()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConn)
	db.SetMaxIdleConns(cfg.db.maxIdleConn)

	// parse a string duration.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}

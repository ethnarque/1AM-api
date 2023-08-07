package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmlogist/1am-audio-srv/internal/data"
	"github.com/pmlogist/1am-audio-srv/internal/mailer"
)

const version = "0.1.0"

type config struct {
	appName string
	port    int
	env     string
	db      struct {
		dsn          string
		maxOpenConns int
		maxIdleTime  string
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
	logger *log.Logger
	models data.Model
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 4000, "API Server Port")
	flag.StringVar(&config.env, "env", os.Getenv("ENV"), "Environment (development|staging|production)")
	flag.StringVar(&config.appName, "app-name", os.Getenv("APP_NAME"), "API Name")

	flag.StringVar(&config.db.dsn, "db-dsn", os.Getenv("DB_URL"), "PostgreSQL DSN")

	flag.IntVar(&config.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.StringVar(&config.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&config.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&config.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&config.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&config.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&config.smtp.port, "smtp-port", 465, "SMTP port")
	flag.StringVar(&config.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&config.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&config.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

	flag.Parse()

	logger := log.New(os.Stdout, config.appName+" | ", log.LstdFlags)

	db, err := openDB(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("database connection pool established")

	mailer, err := mailer.New(
		config.smtp.host,
		config.smtp.port,
		config.smtp.username,
		config.smtp.password,
		config.smtp.sender,
	)
	if err != nil {
		logger.Fatalf("failed to create mail client: %s", err)
	}

	app := &application{
		config: config,
		logger: logger,
		models: data.New(db),
		mailer: mailer,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	dbconfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	dbconfig.MaxConns = int32(cfg.db.maxOpenConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	dbconfig.MaxConnIdleTime = duration

	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

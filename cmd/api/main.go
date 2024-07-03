package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	models "github.com/iamgak/go-ecommerce/internals"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Addr string
}

type Application struct {
	Config          config
	InfoLog         *log.Logger
	ErrorLog        *log.Logger
	DB              *sql.DB
	Cart            models.CartModelInterface
	Order           models.OrderModelInterface
	User            models.UserModelInterface
	Product         models.ProductModelInterface
	Payment         models.PaymentModelInterface
	Seller          models.SellerModelInterface
	Uid             int
	isAuthenticated bool
}

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

func main() {
	var cfg config
	err := godotenv.Load()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbUser := os.Getenv("DB_USERNAME")
	// port := os.Getenv("PORT")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	// addr := flag.String("addr", ":"+port, "HTTP network address")
	dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dsn, "PostgreSQL DSN")
	// Read the connection pool settings from command-line flags into the config struct.
	// Notice the default values that we're using?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rep", 1, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 1, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &Application{
		Config:   cfg,
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		DB:       db,
		User:     &models.UserDB{DB: db},
		Cart:     &models.CartDB{DB: db},
		Order:    &models.OrderDB{DB: db},
		Product:  &models.ProductDB{DB: db},
		Payment:  &models.PaymentDB{DB: db},
		Seller:   &models.SellerDB{DB: db},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.ErrorLog.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatal(err)
	}
}

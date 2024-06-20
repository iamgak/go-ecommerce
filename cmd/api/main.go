package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/iamgak/go-ecommerce/internals/cart"
	"github.com/iamgak/go-ecommerce/internals/checkout"
	"github.com/iamgak/go-ecommerce/internals/object"
	"github.com/iamgak/go-ecommerce/internals/user"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Addr string
}

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	DB            *sql.DB
	Cart          *cart.CartDB
	Checkout      *checkout.OrderDB
	User          *user.UserDB
	Object        *object.ObjectDB
	Authenticated bool
	Uid           int
	// session         *sessions.CookieStore
	// isAuthenticated bool
}

func main() {
	err := godotenv.Load()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbUser := os.Getenv("DB_USERNAME")
	port := os.Getenv("PORT")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	addr := flag.String("addr", ":"+port, "HTTP network address")
	// connStr := flag.String("dsn", fmt.Sprintf("postgres://%s:%s@%s/%s", dbUser, dbPassword, dbHost, dbName), "POSTgres data source name")
	// fmt.Print(*connStr)
	flag.Parse()

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		DB:       db,
		User:     &user.UserDB{DB: db},
		Cart:     &cart.CartDB{DB: db},
		Checkout: &checkout.OrderDB{DB: db},
		Object:   &object.ObjectDB{DB: db},
		// session:  store,
	}
	serve := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	serve.ListenAndServe()
}

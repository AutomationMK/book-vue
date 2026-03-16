package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AutomationMK/book-vue/internal/driver"
)

type config struct {
	port int
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	db       *driver.DB
}

func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatalf("Connot connect to database: %s", err)
	}

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		db:       db,
	}

	err = app.serv()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serv() error {
	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}

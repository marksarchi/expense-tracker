package main

import (
	"context"

	"log"
	"net/http"
	"os"
	"time"

	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/app/handlers"
	"github.com/sarchimark/expense-tracker/foundation/database"
)

func main() {

	log := log.New(os.Stdout, "Expense-Tracker", log.Ldate|log.Ltime|log.Lshortfile)

	//Start expense-tracker app

	if err := Run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}

}

func Run(log *log.Logger) error {

	// dbConfig := database.Config{
	// 	User:       "postgres",
	// 	Password:   "postgres",
	// 	Host:       "localhost:5432",
	// 	Name:       "expensetrackerdb",
	// 	DisableTLS: true,
	// }
	dbConfig := database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       "db",
		Name:       "postgres",
		DisableTLS: true,
	}
	db, err := database.Open(dbConfig)

	if err != nil {
		log.Print(err)
		log.Println("error initializing db")
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := http.Server{
		Addr:    ":8000",
		Handler: handlers.SetupRoutes(db, log, shutdown),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("API listening ")
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil

}

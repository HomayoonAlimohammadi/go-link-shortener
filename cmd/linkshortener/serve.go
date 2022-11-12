package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/homayoonalimohammadi/go-link-shortener/linkshortener/internal/app/database"
	"github.com/spf13/cobra"
)

type LinkShortenerBackend struct {
	Postgres Database
}

const (
	webPort = "8000"
)

var (
	linkShortener = &LinkShortenerBackend{}
	serveCmd      = &cobra.Command{
		Use:   "serve",
		Short: "Serves the Link Shortener",
		Run: func(cmd *cobra.Command, args []string) {
			serve(cmd, args)
		},
	}
)

func serve(cmd *cobra.Command, args []string) {

	// connect to the postgres database
	postgresProvider, err := database.NewPostgresProvider(config.Postgres)
	if err != nil {
		log.Fatalf("could not initialize connection to postgres database: %s", err)
	}
	defer postgresProvider.Close()

	// do the migrations
	err = postgresProvider.Migrate()
	if err != nil {
		log.Println(err)
	}

	linkShortener.Postgres = postgresProvider

	// setup the router
	router := mux.NewRouter()
	router.HandleFunc("/", getRoot)
	router.HandleFunc("/create", createToken).Methods("POST")
	router.HandleFunc("/{token}", redirectToOriginal)
	router.HandleFunc("/{token}/stats", getTokenStats)

	err = http.ListenAndServe(fmt.Sprintf(":%s", webPort), router)

	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

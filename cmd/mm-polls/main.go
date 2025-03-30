package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"log"
	"mm-polls/internal/http-server/cleanup"
	"mm-polls/internal/http-server/create"
	"mm-polls/internal/http-server/finish"
	"mm-polls/internal/http-server/results"
	"mm-polls/internal/http-server/vote"
	"mm-polls/internal/storage/tarantool-db"
	"net/http"
	"os"
)

func main() {
	dbHost := os.Getenv("TARANTOOL_HOST")
	dbPort := os.Getenv("TARANTOOL_PORT")
	dbUser := os.Getenv("TARANTOOL_USER_NAME")
	dbAddr := dbHost + ":" + dbPort

	tt, err := tarantool_db.NewTarantool(dbAddr, dbUser)
	if err != nil {
		log.Fatal("Could not connect to Tarantool ", err)
		return
	}

	err = tt.InitDB()
	if err != nil {
		log.Fatal("Could not initialize Tarantool: ", err)
		return
	}

	log.Println("Tarantool started!")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Post("/create", create.New(tt))
	router.Post("/vote", vote.New(tt))
	router.Post("/results", results.New(tt))
	router.Post("/finish", finish.New(tt))
	router.Post("/cleanup", cleanup.New(tt))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(srv.ListenAndServe())
}

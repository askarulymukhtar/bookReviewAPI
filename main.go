package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	err := openConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer closeConnection()
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Get("/books", handleGetAll)
	v1Router.Get("/books/{id}", handleGetById)
	v1Router.Post("/books", handleCreate)
	v1Router.Put("/books/{id}", handleUpdate)
	v1Router.Delete("/books/{id}", handleDelete)
	router.Mount("/v1", v1Router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not identified")
	}

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("starting server on port: %v", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/askarulymukhtar/go/bookReviewAPI/db"
	"github.com/askarulymukhtar/go/bookReviewAPI/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	err := db.OpenConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConnection()
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
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerError)
	v1Router.Get("/books", handlers.HandleGetAll)
	v1Router.Get("/books/{id}", handlers.HandleGetById)
	v1Router.Post("/books", handlers.HandleCreate)
	v1Router.Put("/books/{id}", handlers.HandleUpdate)
	v1Router.Delete("/books/{id}", handlers.HandleDelete)
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

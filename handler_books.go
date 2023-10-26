package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type BookDTO struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
	Genre           string `json:"genre"`
	ISNB            string `json:"isnb"`
}

func handleGetById(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id := queryParams.Get("id")

	var book BookDTO

	if err := DB.QueryRow("SELECT title, author, publication_year, genre, isnb FROM books WHERE id = $1", id).Scan(
		&book.Title, &book.Author, &book.PublicationYear, &book.Genre, &book.ISNB); err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	params := BookDTO{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("error marshalling json")
		return
	}

	if err := DB.QueryRow("INSERT INTO books(id, title, author, publication_year, genre, isnb) VALUES ($1, $2, $3, $4, $5, $6)",
		uuid.New(), params.Title, params.Author, params.PublicationYear, params.Genre, params.ISNB).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

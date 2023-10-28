package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func handleGetAll(w http.ResponseWriter, r *http.Request) {

	rows, err := DB.Query("SELECT title, author, publication_year, genre, isnb FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var books []BookDTO
	for rows.Next() {
		var i BookDTO
		if err = rows.Scan(&i.Title, &i.Author, &i.PublicationYear, &i.Genre, &i.ISNB); err != nil {
			log.Fatal(err)
		}
		books = append(books, i)
	}
	json.NewEncoder(w).Encode(books)
}

func handleGetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	params := BookDTO{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("error marshalling json")
		return
	}

	if err := DB.QueryRow("UPDATE books SET title=$2, author=$3, publication_year=$4, genre=$5, isnb=$6 WHERE id=$1",
		id, params.Title, params.Author, params.PublicationYear, params.Genre, params.ISNB).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := DB.QueryRow("DELETE FROM books WHERE id=$1", id).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

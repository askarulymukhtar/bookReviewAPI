package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/askarulymukhtar/go/bookReviewAPI/db"
	"github.com/askarulymukhtar/go/bookReviewAPI/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func HandleGetAll(w http.ResponseWriter, r *http.Request) {

	rows, err := db.DB.Query("SELECT title, author, publication_year, genre, isnb FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var books []models.BookDTO
	for rows.Next() {
		var i models.BookDTO
		if err = rows.Scan(&i.Title, &i.Author, &i.PublicationYear, &i.Genre, &i.ISNB); err != nil {
			log.Fatal(err)
		}
		books = append(books, i)
	}
	json.NewEncoder(w).Encode(books)
}

func HandleGetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var book models.BookDTO

	if err := db.DB.QueryRow("SELECT title, author, publication_year, genre, isnb FROM books WHERE id = $1", id).Scan(
		&book.Title, &book.Author, &book.PublicationYear, &book.Genre, &book.ISNB); err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	params := models.BookDTO{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("error marshalling json")
		return
	}

	if err := db.DB.QueryRow("INSERT INTO books(id, title, author, publication_year, genre, isnb) VALUES ($1, $2, $3, $4, $5, $6)",
		uuid.New(), params.Title, params.Author, params.PublicationYear, params.Genre, params.ISNB).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	params := models.BookDTO{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("error marshalling json")
		return
	}

	if err := db.DB.QueryRow("UPDATE books SET title=$2, author=$3, publication_year=$4, genre=$5, isnb=$6 WHERE id=$1",
		id, params.Title, params.Author, params.PublicationYear, params.Genre, params.ISNB).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := db.DB.QueryRow("DELETE FROM books WHERE id=$1", id).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

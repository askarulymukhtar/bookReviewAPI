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
	type result struct {
		books []models.BookDTO
		err   error
	}

	// Create a channel to communicate the results from the goroutine
	resultChannel := make(chan result)

	go func() {
		rows, err := db.DB.Query("SELECT title, author, publication_year, genre, isbn FROM books")
		if err != nil {
			resultChannel <- result{err: err}
			return
		}
		defer rows.Close()

		var books []models.BookDTO
		for rows.Next() {
			var i models.BookDTO
			if err = rows.Scan(&i.Title, &i.Author, &i.PublicationYear, &i.Genre, &i.ISBN); err != nil {
				resultChannel <- result{err: err}
				return
			}
			books = append(books, i)
		}
		resultChannel <- result{books: books}
	}()

	res := <-resultChannel
	if res.err != nil {
		log.Printf("error while getting books: %v", res.err)
		http.Error(w, res.err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.books)
}

func HandleGetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Define a result struct to hold the query results
	type result struct {
		book models.BookDTO
		err  error
	}

	// Create a channel to communicate results
	resultChannel := make(chan result)

	// Execute the query in a goroutine
	go func() {
		var book models.BookDTO
		err := db.DB.QueryRow("SELECT title, author, publication_year, genre, isbn FROM books WHERE id = $1", id).Scan(
			&book.Title, &book.Author, &book.PublicationYear, &book.Genre, &book.ISBN)

		// Send the results back to the main goroutine
		resultChannel <- result{book: book, err: err}
	}()

	// Receive the results
	res := <-resultChannel

	// Handle the error if any
	if res.err != nil {
		log.Printf("error while getting book: %v", res.err)
		http.Error(w, "could not retrieve book", http.StatusInternalServerError)
		return
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.book)
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	params := models.BookDTO{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while parsing json: %v", err)
		return
	}

	resultChannel := make(chan error)

	go func() {
		err := db.DB.QueryRow("INSERT INTO books(id, title, author, publication_year, genre, isbn) VALUES ($1, $2, $3, $4, $5, $6)",
			uuid.New(), params.Title, params.Author, params.PublicationYear, params.Genre, params.ISBN).Err()
		resultChannel <- err
	}()

	err := <-resultChannel
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while creating book: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	params := models.BookDTO{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while parsing json: %v", err)
		return
	}

	resultChannel := make(chan error)

	go func() {
		err := db.DB.QueryRow("UPDATE books SET title=$2, author=$3, publication_year=$4, genre=$5, isbn=$6 WHERE id=$1",
			id, params.Title, params.Author, params.PublicationYear, params.Genre, params.ISBN).Err()
		resultChannel <- err
	}()

	err := <-resultChannel
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while updating book: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resultChannel := make(chan error)

	go func() {
		err := db.DB.QueryRow("DELETE FROM books WHERE id=$1", id).Err()
		resultChannel <- err
	}()

	err := <-resultChannel
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while deleting book: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

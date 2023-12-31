package models

import (
	"github.com/google/uuid"
)

type Book struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	PublicationYear int       `json:"publication_year"`
	Genre           string    `json:"genre"`
	ISBN            string    `json:"isbn"`
}

type BookDTO struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
	Genre           string `json:"genre"`
	ISBN            string `json:"isbn"`
}

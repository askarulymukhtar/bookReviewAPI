package main

import (
	"github.com/google/uuid"
)

type Book struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	PublicationYear int       `json:"publication_year"`
	Genre           string    `json:"genre"`
	ISNB            string    `json:"isnb"`
}

type BookDTO struct {
	Title           string
	Author          string
	PublicationYear int
	Genre           string
	ISBN            string
}

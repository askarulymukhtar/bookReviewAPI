-- name: CreateBook :one
INSERT INTO books (id, title, author, publication_year, genre, isnb)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
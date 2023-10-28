-- +goose Up
CREATE TABLE books (
    id UUID UNIQUE NOT NULL PRIMARY KEY,
    title TEXT,
    author TEXT,
    publication_year INTEGER,
    genre TEXT,
    isbn TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE books;
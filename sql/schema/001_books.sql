-- +goose Up
CREATE TABLE books (
    id UUID PRIMARY KEY,
    title TEXT,
    author TEXT,
    publication_year INTEGER,
    genre TEXT,
    isnb TEXT
);

-- +goose Down
DROP TABLE books;
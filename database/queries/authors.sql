-- name: AuthorsGetAll :many
SELECT * FROM authors;

-- name: AuthorsInsert :exec
INSERT INTO authors (name, email, created_at)
VALUES ($2, $1, NOW());

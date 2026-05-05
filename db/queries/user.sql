-- name: Create :one
INSERT INTO users (email, firstname, lastname, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindByEmail :one
SELECT * FROM users WHERE email = $1;


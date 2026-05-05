-- name: Create :one
INSERT INTO users (email, firstname, lastname, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: FindByID :one
SELECT * FROM users WHERE ID = $1;


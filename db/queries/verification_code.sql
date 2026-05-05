-- name: CreateVerificationCode :exec
INSERT INTO verification_codes (user_id, code, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVerificationCode :one
SELECT * FROM verification_codes
WHERE user_id = $1;

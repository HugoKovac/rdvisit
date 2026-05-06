-- name: GetAllCenters :many
SELECT * FROM centers;

-- name: GetCenterByID :one
SELECT * FROM centers
WHERE id = $1;

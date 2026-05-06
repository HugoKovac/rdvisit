-- name: CreateAppointment :one
INSERT INTO appointments (date, practitioner_id, patient_email)
VALUES ($1, $2, $3)
RETURNING *;

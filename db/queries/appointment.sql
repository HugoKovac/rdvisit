-- name: CreateAppointment :one
WITH patient AS (
    INSERT INTO users (email, firstname, lastname, password_hash, role, verified)
    VALUES ($3, '', '', '', 'patient', false)
    ON CONFLICT (email) DO UPDATE
    SET email = EXCLUDED.email
    RETURNING id, email, role
)
INSERT INTO appointments (date, practitioner_id, patient_id)
SELECT $1, $2, patient.id
FROM patient
WHERE patient.role = 'patient'
RETURNING *;

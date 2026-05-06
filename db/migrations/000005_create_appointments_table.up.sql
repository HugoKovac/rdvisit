CREATE TABLE appointments (
    id UUID             PRIMARY KEY DEFAULT gen_random_uuid(),

    date                TIMESTAMPTZ NOT NULL,

    practitioner_id     UUID NOT NULL REFERENCES users(id),
    patient_email       TEXT NOT NULL,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE centers (
    id UUID             PRIMARY KEY DEFAULT gen_random_uuid(),

    name                TEXT NOT NULL,
    street              TEXT NOT NULL,
    street_number       TEXT NOT NULL,
    city                TEXT NOT NULL,
    postal_code         TEXT NOT NULL,
    region              TEXT NOT NULL,
    country             TEXT NOT NULL,

    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

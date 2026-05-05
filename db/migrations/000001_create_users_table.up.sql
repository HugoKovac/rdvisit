CREATE TABLE IF NOT EXISTS users (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT        NOT NULL UNIQUE,
    firstname       TEXT        NOT NULL,
    lastname        TEXT        NOT NULL,
    password_hash   TEXT        NOT NULL,
    verified        BOOLEAN     DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
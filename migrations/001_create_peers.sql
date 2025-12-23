CREATE TABLE peers (
    id SERIAL PRIMARY KEY,

    public_key TEXT NOT NULL UNIQUE,
    vpn_address TEXT NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
BEGIN;

DROP TABLE IF EXISTS peers;

CREATE TABLE peers (
    id SERIAL PRIMARY KEY,
    common_name TEXT NOT NULL UNIQUE,
    node_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE peer_keys (
    id SERIAL PRIMARY KEY,
    peer_id INTEGER NOT NULL REFERENCES peers(id) ON DELETE CASCADE,
    cert_pem TEXT NOT NULL,
    key_pem TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;
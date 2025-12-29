DROP TABLE IF EXISTS peer_keys;

DROP TABLE IF EXISTS peers;

CREATE TABLE peers (
    id SERIAL PRIMARY KEY,
    common_name TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
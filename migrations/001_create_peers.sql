DROP TABLE IF EXISTS peer_keys;
DROP TABLE IF EXISTS peers;

DROP TYPE IF EXISTS peer_connection_status;

CREATE TYPE peer_connection_status AS ENUM (
    'connected',
    'not_connected'
);

CREATE TABLE peers (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL UNIQUE,
    telegram_id BIGINT NOT NULL,
    connection_status peer_connection_status NOT NULL DEFAULT 'not_connected',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_peers_telegram_id ON peers (telegram_id);

ALTER TABLE peers
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

UPDATE peers SET is_active = TRUE WHERE is_active IS NULL;
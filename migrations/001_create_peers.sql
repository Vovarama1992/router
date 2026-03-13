CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    active_until TIMESTAMP NOT NULL DEFAULT now() + interval '30 days',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE peers (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL UNIQUE,
    telegram_id BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now()
);
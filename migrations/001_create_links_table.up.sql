CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_slug ON urls(slug);
CREATE INDEX idx_url ON urls(url);

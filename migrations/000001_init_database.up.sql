CREATE TABLE slugs (
    id SERIAL PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    next_number INT DEFAULT(1) NOT NULL
);

CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    slug_id INT REFERENCES slugs(id) NOT NULL,
    slug_number INT NOT NULL
);

CREATE UNIQUE INDEX idx_unique_url_slug ON urls (slug_id, slug_number);

ALTER TABLE urls ADD slug_id SERIAL;
CREATE UNIQUE INDEX idx_urls_slug_slug_id ON urls (slug, slug_id);

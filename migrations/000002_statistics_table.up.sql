CREATE TABLE redirects (
    id BIGSERIAL PRIMARY KEY,
    url_id INTEGER NOT NULL REFERENCES urls(id),
    redirected_at TIMESTAMP NOT NULL,
    ip VARCHAR(15) NOT NULL,
    user_agent TEXT
)

CREATE INDEX idx_redirects_url_id ON redirects(url_id);

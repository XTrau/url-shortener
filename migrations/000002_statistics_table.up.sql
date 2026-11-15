CREATE TABLE redirects (
    id BIGSERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id) NOT NULL,
    redirected_at TIMESTAMP NOT NULL,
    ip VARCHAR(15) NOT NULL,
    user_agent TEXT
);

CREATE INDEX idx_redirects_url_id ON redirects USING HASH(url_id);
CREATE INDEX idx_redirects_redirected_at ON redirects USING BTREE(url_id);

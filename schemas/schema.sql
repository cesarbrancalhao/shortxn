CREATE TABLE urls (
    id VARCHAR(255) PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    clicks BIGINT DEFAULT 0,
    UNIQUE(long_url),
    UNIQUE(short_url)
);

CREATE INDEX idx_long_url ON urls(long_url);
CREATE INDEX idx_short_url ON urls(short_url);

CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    url_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(45),
    FOREIGN KEY (url_id) REFERENCES urls(id)
);

CREATE INDEX idx_analytics_url_id ON analytics(url_id);
CREATE INDEX idx_analytics_timestamp ON analytics(timestamp);

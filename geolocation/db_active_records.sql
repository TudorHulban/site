CREATE TABLE ip_logs (
    id BIGSERIAL PRIMARY KEY,
    ip_address VARCHAR(45) NOT NULL,
    postcode VARCHAR(20),
    city_id INT REFERENCES cities(id),
    asn_id INT REFERENCES asn_numbers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_raw_hits_time ON raw_ip_hits(created_at);
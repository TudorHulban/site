CREATE TABLE aggregates_city (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    interval_start DATETIME NOT NULL,
    interval_end DATETIME NOT NULL,
    city_id INTEGER NOT NULL,
    hit_count INTEGER NOT NULL,
    UNIQUE(interval_start, interval_end, city_id)
);

CREATE INDEX idx_city_interval ON aggregates_city(interval_start, interval_end);


CREATE TABLE aggregates_asn (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    interval_start DATETIME NOT NULL,
    interval_end DATETIME NOT NULL,
    asn_id INTEGER NOT NULL,
    hit_count INTEGER NOT NULL,
    UNIQUE(interval_start, interval_end, asn_id)
);

CREATE INDEX idx_asn_interval ON aggregates_asn(interval_start, interval_end);
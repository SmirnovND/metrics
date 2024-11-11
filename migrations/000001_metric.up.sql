CREATE TABLE metric (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION,
    delta INT
);

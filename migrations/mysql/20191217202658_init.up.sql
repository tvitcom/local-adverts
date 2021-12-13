CREATE TABLE album
(
    id         VARCHAR(128) PRIMARY KEY,
    name       VARCHAR(128) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

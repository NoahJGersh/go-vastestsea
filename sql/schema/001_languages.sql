-- +goose Up
CREATE TABLE languages (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL
);

CREATE INDEX ON languages (LOWER(name));

-- +goose Down
DROP TABLE languages;
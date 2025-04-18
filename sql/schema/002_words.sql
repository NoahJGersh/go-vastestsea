-- +goose Up
CREATE TABLE words (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    word TEXT NOT NULL,
    font_formatted TEXT, -- Nullable, because some words require alternate formatting
    language_id UUID NOT NULL,
    CONSTRAINT fk_language_id
    FOREIGN KEY (language_id)
    REFERENCES languages(id)
    ON DELETE CASCADE,
    UNIQUE (language_id, word)
);

-- +goose Down
DROP TABLE words;
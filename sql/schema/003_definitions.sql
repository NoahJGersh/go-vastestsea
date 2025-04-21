-- +goose Up
CREATE TABLE definitions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    content TEXT NOT NULL,
    part_of_speech TEXT NOT NULL,
    word_id UUID NOT NULL,
    CONSTRAINT fk_word_id
    FOREIGN KEY (word_id)
    REFERENCES words(id)
    ON DELETE CASCADE,
    UNIQUE (word_id, content)
);

-- name: CreateWord :one
INSERT INTO words (id, created_at, updated_at, word, language_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: CreateFormattedWord :one
INSERT INTO words (id, created_at, updated_at, word, font_formatted, language_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetWord :many
SELECT * FROM words
WHERE LOWER(word) = $1;

-- name: GetWordByID :one
SELECT * FROM words
WHERE id = $1;

-- name: GetWordFromLanguage :one
SELECT * FROM words
WHERE LOWER(word) = $1 AND language_id = $2;

-- name: GetWords :many
SELECT * FROM words;

-- name: GetWordsByLanguageID :many
SELECT * FROM words
WHERE language_id = $1;

-- name: UpdateWord :one
UPDATE words
SET 
    word = CASE WHEN @set_word::bool
        THEN @word::text
        ELSE word
        END,
    font_formatted = CASE WHEN @set_formatted::bool
        THEN @formatted::text
        ELSE font_formatted
        END
WHERE id = @id
RETURNING *;

-- name: DeleteWord :exec
DELETE FROM words
WHERE id = $1;
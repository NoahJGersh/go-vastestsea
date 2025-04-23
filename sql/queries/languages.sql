-- name: CreateLanguage :one
INSERT INTO languages (id, created_at, updated_at, name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetLanguages :many
SELECT * FROM languages;

-- name: GetLanguage :one
SELECT * FROM languages
WHERE LOWER(name) = $1;

-- name: GetLanguageByID :one
SELECT * FROM languages
WHERE id = $1;

-- name: DeleteLanguage :exec
DELETE FROM languages
WHERE id = $1;

-- name: UpdateLanguageName :one
UPDATE languages
SET name = $1
WHERE LOWER(name) = $2
RETURNING *;
-- name: CreateLanguage :one
INSERT INTO languages (id, created_at, updated_at, name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetLanguage :one
SELECT * FROM languages
WHERE LOWER(name) = $1;
-- name: CreateDefinition :one
INSERT INTO definitions (id, created_at, updated_at, content, part_of_speech, word_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetDefinitionsOfWord :many
SELECT * FROM definitions
WHERE definitions.word_id = $1
ORDER BY part_of_speech ASC, content ASC;

-- name: GetDefinitionByID :one
SELECT * FROM definitions
WHERE id = $1;

-- name: UpdateDefinitionPartOfSpeech :one
UPDATE definitions
SET part_of_speech = $1
WHERE id = $2
RETURNING *;

-- name: UpdateDefinitionContent :one
UPDATE definitions
SET content = $1
WHERE id = $2
RETURNING *;

-- name: UpdateDefinition :one
UPDATE definitions
SET content = $1, part_of_speech = $2
WHERE id = $3
RETURNING *;

-- name: DeleteDefinition :exec
DELETE FROM definitions
WHERE id = $1;
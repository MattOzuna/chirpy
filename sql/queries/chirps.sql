-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid (),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirpsAsc :many
SELECT * FROM chirps 
ORDER BY created_at ASC;

-- name: GetAllChirpsDesc :many
SELECT * FROM chirps 
ORDER BY created_at DeSC;

-- name: GetAllChirpsByUserIDAsc :many
SELECT * FROM chirps 
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetAllChirpsByUserIDDesc :many
SELECT * FROM chirps 
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteAllChirp :exec
DELETE FROM chirps;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1; 
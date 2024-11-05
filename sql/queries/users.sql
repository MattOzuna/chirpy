-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid (),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: EditUserPassword :one
UPDATE users
SET hashed_password = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: EditUserEmail :one
UPDATE users
SET email = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: EditUser :one
UPDATE users
SET hashed_password = $1, email = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;
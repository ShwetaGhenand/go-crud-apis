-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  id, name, email, phone, age, address
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    name = $1,
	email = $2,
	phone = $3,
	age = $4,
	address = $5
	WHERE id = $6
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

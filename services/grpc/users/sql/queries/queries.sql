-- name: UserExists :one
SELECT * FROM users
WHERE name = $1 and password = $2;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (
  id, name, password, email, phone, age, address, createTime, updateTime
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: UpdateUser :exec
UPDATE users SET
    name = $1,
	email = $2,
	phone = $3,
	age = $4,
	address = $5,
	updateTime = $6
	WHERE id = $7;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

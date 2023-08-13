-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1  LIMIT 1;

-- name: UpdateUser :one
update 
  users
set 
  hashed_password = coalesce(sqlc.narg(hashed_password), hashed_password),
  full_name = coalesce(sqlc.narg(full_name), full_name),
  email = coalesce(sqlc.narg(email), email)
where 
  username = sqlc.arg(username)
returning *;

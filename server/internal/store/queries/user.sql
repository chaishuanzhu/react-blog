-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :execresult
INSERT INTO users (email, password_hash, nickname, avatar, website, role)
VALUES (?, ?, ?, ?, ?, 'admin');

-- name: UpdateUserPassword :execrows
UPDATE users SET password_hash = ?, token_version = token_version + 1 WHERE id = ?;

-- name: UpdateUserProfile :execrows
UPDATE users SET nickname = ?, avatar = ?, website = ? WHERE id = ?;

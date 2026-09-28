-- name: GetPublishedArticleRef :one
SELECT id, title FROM articles WHERE id = ? AND status = 'published';

-- name: GetArticleRef :one
SELECT id, title FROM articles WHERE id = ?;

-- name: ListTopLevelComments :many
SELECT id, article_id, parent_id, nickname, website, avatar, content, is_admin, created_at
FROM comments
WHERE article_id <=> sqlc.narg(article_id) AND parent_id IS NULL
ORDER BY created_at DESC, id DESC
LIMIT ? OFFSET ?;

-- name: CountTopLevelComments :one
SELECT COUNT(*) FROM comments
WHERE article_id <=> sqlc.narg(article_id) AND parent_id IS NULL;

-- name: ListRepliesForParents :many
SELECT id, article_id, parent_id, nickname, website, avatar, content, is_admin, created_at
FROM comments
WHERE parent_id IN (sqlc.slice(parent_ids))
ORDER BY created_at ASC, id ASC;

-- name: GetComment :one
SELECT id, article_id, parent_id, nickname, email, content, is_admin
FROM comments WHERE id = ?;

-- name: CreateComment :execresult
INSERT INTO comments (article_id, parent_id, nickname, email, website, avatar, content, is_admin, ip, user_agent)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetCommentByID :one
SELECT id, article_id, parent_id, nickname, website, avatar, content, is_admin, created_at
FROM comments WHERE id = ?;

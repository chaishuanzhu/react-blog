-- name: CreateCategory :execresult
INSERT INTO categories (name) VALUES (?);

-- name: RenameCategory :execrows
UPDATE categories SET name = ? WHERE id = ?;

-- name: DeleteCategory :execrows
DELETE FROM categories WHERE id = ?;

-- name: CreateTag :execresult
INSERT INTO tags (name) VALUES (?);

-- name: RenameTag :execrows
UPDATE tags SET name = ? WHERE id = ?;

-- name: DeleteTag :execrows
DELETE FROM tags WHERE id = ?;

-- name: AdminListComments :many
SELECT cm.id, cm.article_id, cm.parent_id, cm.nickname, cm.email, cm.website, cm.avatar,
       cm.content, cm.is_admin, cm.ip, cm.created_at,
       a.title AS article_title
FROM comments cm
LEFT JOIN articles a ON a.id = cm.article_id
ORDER BY cm.created_at DESC, cm.id DESC
LIMIT ? OFFSET ?;

-- name: AdminCountComments :one
SELECT COUNT(*) FROM comments;

-- name: DeleteComment :execrows
DELETE FROM comments WHERE id = ?;

-- name: CreateMoment :execresult
INSERT INTO moments (content, images, created_at) VALUES (?, ?, ?);

-- name: UpdateMoment :execrows
UPDATE moments SET content = ?, images = ?, created_at = ? WHERE id = ?;

-- name: DeleteMoment :execrows
DELETE FROM moments WHERE id = ?;

-- name: CreateFriendLink :execresult
INSERT INTO friend_links (name, url, avatar, description) VALUES (?, ?, ?, ?);

-- name: UpdateFriendLink :execrows
UPDATE friend_links SET name = ?, url = ?, avatar = ?, description = ? WHERE id = ?;

-- name: DeleteFriendLink :execrows
DELETE FROM friend_links WHERE id = ?;

-- name: CreateChangelog :execresult
INSERT INTO changelogs (items, logged_at) VALUES (?, ?);

-- name: UpdateChangelog :execrows
UPDATE changelogs SET items = ?, logged_at = ? WHERE id = ?;

-- name: DeleteChangelog :execrows
DELETE FROM changelogs WHERE id = ?;

-- name: UpdatePage :execrows
UPDATE pages SET content = ? WHERE page_key = ?;

-- name: AdminStats :one
SELECT
  (SELECT COUNT(*) FROM articles WHERE status = 'published') AS published_count,
  (SELECT COUNT(*) FROM articles WHERE status = 'draft') AS draft_count,
  (SELECT COUNT(*) FROM categories) AS category_count,
  (SELECT COUNT(*) FROM tags) AS tag_count,
  (SELECT COUNT(*) FROM comments) AS comment_count,
  (SELECT COUNT(*) FROM moments) AS moment_count,
  (SELECT COUNT(*) FROM friend_links) AS friend_link_count,
  (SELECT view_count FROM site_settings WHERE id = 1) AS view_count;

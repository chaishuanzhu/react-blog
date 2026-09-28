-- name: ListMoments :many
SELECT id, content, images, created_at FROM moments
ORDER BY created_at DESC, id DESC
LIMIT ? OFFSET ?;

-- name: CountMoments :one
SELECT COUNT(*) FROM moments;

-- name: ListFriendLinks :many
SELECT id, name, url, avatar, description, created_at FROM friend_links
ORDER BY created_at DESC, id DESC;

-- name: ListChangelogs :many
SELECT id, items, logged_at FROM changelogs
ORDER BY logged_at DESC, id DESC;

-- name: ListProjects :many
SELECT id, name, description, cover, url, sort_order FROM projects
ORDER BY sort_order ASC, id ASC;

-- name: GetPage :one
SELECT page_key, content, updated_at FROM pages WHERE page_key = ?;

-- name: GetSiteSettings :one
SELECT id, notice, view_count, updated_at FROM site_settings WHERE id = 1;

-- LAST_INSERT_ID(expr) makes the new value available via Result.LastInsertId without a second query.
-- name: IncrementViewCount :execresult
UPDATE site_settings SET view_count = LAST_INSERT_ID(view_count + 1) WHERE id = 1;

-- name: UpdateNotice :exec
UPDATE site_settings SET notice = ? WHERE id = 1;

-- name: GetSiteStats :one
SELECT
  (SELECT COUNT(*) FROM articles WHERE status = 'published') AS article_count,
  (SELECT COUNT(*) FROM categories) AS category_count,
  (SELECT COUNT(*) FROM tags) AS tag_count;

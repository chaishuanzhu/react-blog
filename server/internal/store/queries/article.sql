-- name: ListPublishedArticles :many
SELECT a.id, a.title, a.content, a.published_at, a.updated_at,
       a.category_id, c.name AS category_name
FROM articles a
LEFT JOIN categories c ON c.id = a.category_id
WHERE a.status = 'published'
  AND (sqlc.narg(keyword) IS NULL OR a.title LIKE sqlc.narg(keyword))
  AND (sqlc.narg(category) IS NULL OR c.name = sqlc.narg(category))
  AND (sqlc.narg(tag) IS NULL OR EXISTS (
        SELECT 1 FROM article_tags at JOIN tags t ON t.id = at.tag_id
        WHERE at.article_id = a.id AND t.name = sqlc.narg(tag)))
ORDER BY a.published_at DESC, a.id DESC
LIMIT ? OFFSET ?;

-- name: CountPublishedArticles :one
SELECT COUNT(*)
FROM articles a
LEFT JOIN categories c ON c.id = a.category_id
WHERE a.status = 'published'
  AND (sqlc.narg(keyword) IS NULL OR a.title LIKE sqlc.narg(keyword))
  AND (sqlc.narg(category) IS NULL OR c.name = sqlc.narg(category))
  AND (sqlc.narg(tag) IS NULL OR EXISTS (
        SELECT 1 FROM article_tags at JOIN tags t ON t.id = at.tag_id
        WHERE at.article_id = a.id AND t.name = sqlc.narg(tag)));

-- name: GetPublishedArticle :one
SELECT a.id, a.title, a.content, a.published_at, a.updated_at,
       a.category_id, c.name AS category_name
FROM articles a
LEFT JOIN categories c ON c.id = a.category_id
WHERE a.id = ? AND a.status = 'published';

-- name: ListTagsForArticles :many
SELECT at.article_id, t.id, t.name
FROM article_tags at
JOIN tags t ON t.id = at.tag_id
WHERE at.article_id IN (sqlc.slice(article_ids))
ORDER BY t.name;

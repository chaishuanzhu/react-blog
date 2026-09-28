-- name: AdminListArticles :many
SELECT a.id, a.title, a.content, a.status, a.published_at, a.created_at, a.updated_at,
       a.category_id, c.name AS category_name
FROM articles a
LEFT JOIN categories c ON c.id = a.category_id
WHERE (sqlc.narg(status) IS NULL OR a.status = sqlc.narg(status))
  AND (sqlc.narg(keyword) IS NULL OR a.title LIKE sqlc.narg(keyword))
  AND (sqlc.narg(category_id) IS NULL OR a.category_id = sqlc.narg(category_id))
  AND (sqlc.narg(tag_id) IS NULL OR EXISTS (
        SELECT 1 FROM article_tags at WHERE at.article_id = a.id AND at.tag_id = sqlc.narg(tag_id)))
ORDER BY a.published_at DESC, a.id DESC
LIMIT ? OFFSET ?;

-- name: AdminCountArticles :one
SELECT COUNT(*)
FROM articles a
WHERE (sqlc.narg(status) IS NULL OR a.status = sqlc.narg(status))
  AND (sqlc.narg(keyword) IS NULL OR a.title LIKE sqlc.narg(keyword))
  AND (sqlc.narg(category_id) IS NULL OR a.category_id = sqlc.narg(category_id))
  AND (sqlc.narg(tag_id) IS NULL OR EXISTS (
        SELECT 1 FROM article_tags at WHERE at.article_id = a.id AND at.tag_id = sqlc.narg(tag_id)));

-- name: AdminGetArticle :one
SELECT a.id, a.title, a.content, a.status, a.published_at, a.created_at, a.updated_at,
       a.category_id, c.name AS category_name
FROM articles a
LEFT JOIN categories c ON c.id = a.category_id
WHERE a.id = ?;

-- name: CreateArticle :execresult
INSERT INTO articles (title, content, category_id, status, published_at)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateArticle :execrows
UPDATE articles
SET title = ?, content = ?, category_id = ?, status = ?, published_at = ?
WHERE id = ?;

-- name: DeleteArticle :execrows
DELETE FROM articles WHERE id = ?;

-- name: ArticleExists :one
SELECT COUNT(*) FROM articles WHERE id = ?;

-- name: DeleteArticleTags :exec
DELETE FROM article_tags WHERE article_id = ?;

-- name: AddArticleTag :exec
INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?);

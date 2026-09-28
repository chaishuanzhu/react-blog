-- name: ListCategoriesWithCount :many
SELECT c.id, c.name, c.created_at,
       CAST(COUNT(a.id) AS SIGNED) AS article_count
FROM categories c
LEFT JOIN articles a ON a.category_id = c.id AND a.status = 'published'
GROUP BY c.id, c.name, c.created_at
ORDER BY c.created_at DESC, c.id DESC;

-- name: CountUncategorizedPublished :one
SELECT COUNT(*) FROM articles WHERE status = 'published' AND category_id IS NULL;

-- name: ListTagsWithCount :many
SELECT t.id, t.name, t.created_at,
       CAST(COUNT(a.id) AS SIGNED) AS article_count
FROM tags t
LEFT JOIN article_tags at ON at.tag_id = t.id
LEFT JOIN articles a ON a.id = at.article_id AND a.status = 'published'
GROUP BY t.id, t.name, t.created_at
ORDER BY t.created_at DESC, t.id DESC;

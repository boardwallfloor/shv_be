-- name: GetArticle :one
SELECT * FROM articles
WHERE id = ? LIMIT 1;

-- name: ListArticles :many
SELECT * FROM articles
ORDER BY created_date DESC
LIMIT ?
OFFSET ?;

-- name: CreateArticle :execresult
INSERT INTO articles (
  title, content, category, status
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateArticle :exec
UPDATE articles
SET 
  title = ?,
  content = ?,
  category = ?,
  status = ?
WHERE id = ?;

-- name: DeleteArticle :exec
DELETE FROM articles
WHERE id = ?;


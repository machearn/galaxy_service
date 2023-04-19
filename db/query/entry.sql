-- name: GetEntry :one
SELECT * FROM entries WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateEntry :one
INSERT INTO entries (
  member_id, item_id, quantity, total, created_at
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateEntry :one
UPDATE entries SET
  member_id = coalesce(sqlc.narg('member_id'), member_id),
  item_id = coalesce(sqlc.narg('item_id'), item_id),
  quantity = coalesce(sqlc.narg('quantity'), quantity),
  total = coalesce(sqlc.narg('total'), total),
  created_at = coalesce(sqlc.narg('created_at'), created_at)
WHERE id = @id
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id = $1;
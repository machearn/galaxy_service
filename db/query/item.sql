-- name: GetItem :one
SELECT * FROM items WHERE id = $1 LIMIT 1;

-- name: ListItems :many
SELECT * FROM items
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateItem :one
INSERT INTO items (
  name, amount, price
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateItem :one
UPDATE items SET
  name = coalesce(sqlc.narg('name'), name),
  amount = coalesce(sqlc.narg('amount'), amount),
  price = coalesce(sqlc.narg('price'), price)
WHERE id = @id
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1;
// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: item.sql

package db

import (
	"context"
	"database/sql"
)

const createItem = `-- name: CreateItem :one
INSERT INTO items (
  name, amount, price
) VALUES (
  $1, $2, $3
) RETURNING id, name, amount, price
`

type CreateItemParams struct {
	Name   string `json:"name"`
	Amount int32  `json:"amount"`
	Price  int32  `json:"price"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, createItem, arg.Name, arg.Amount, arg.Price)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Amount,
		&i.Price,
	)
	return i, err
}

const deleteItem = `-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1
`

func (q *Queries) DeleteItem(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteItem, id)
	return err
}

const getItem = `-- name: GetItem :one
SELECT id, name, amount, price FROM items WHERE id = $1 LIMIT 1
`

func (q *Queries) GetItem(ctx context.Context, id int32) (Item, error) {
	row := q.db.QueryRowContext(ctx, getItem, id)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Amount,
		&i.Price,
	)
	return i, err
}

const listItems = `-- name: ListItems :many
SELECT id, name, amount, price FROM items
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListItemsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListItems(ctx context.Context, arg ListItemsParams) ([]Item, error) {
	rows, err := q.db.QueryContext(ctx, listItems, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Amount,
			&i.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateItem = `-- name: UpdateItem :one
UPDATE items SET
  name = coalesce($1, name),
  amount = coalesce($2, amount),
  price = coalesce($3, price)
WHERE id = $4
RETURNING id, name, amount, price
`

type UpdateItemParams struct {
	Name   sql.NullString `json:"name"`
	Amount sql.NullInt32  `json:"amount"`
	Price  sql.NullInt32  `json:"price"`
	ID     int32          `json:"id"`
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, updateItem,
		arg.Name,
		arg.Amount,
		arg.Price,
		arg.ID,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Amount,
		&i.Price,
	)
	return i, err
}

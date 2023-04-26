package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	ListEntriesByMemberTx(ctx context.Context, arg ListEntriesByMemberTxParams) ([]Entry, error)
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txQuery := New(tx)
	err = fn(txQuery)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("execute err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

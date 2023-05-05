package db

import (
	"context"
)

type ListEntriesByMemberTxParams struct {
	Username string
	Limit    int32
	Offset   int32
}

func (store *SQLStore) ListEntriesByMemberTx(ctx context.Context, arg ListEntriesByMemberTxParams) ([]Entry, error) {
	var result []Entry

	err := store.execTx(ctx, func(q *Queries) error {
		member, err := q.GetMemberByName(ctx, arg.Username)
		if err != nil {
			return err
		}

		arg := ListEntriesByMemberParams{
			MemberID: member.ID,
			Limit:    arg.Limit,
			Offset:   arg.Offset,
		}
		result, err = q.ListEntriesByMember(ctx, arg)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (store *SQLStore) CreateMemberTx(ctx context.Context, arg CreateMemberParams, callback func(Member) error) (Member, error) {
	var result Member

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result, err = q.CreateMember(ctx, arg)
		if err != nil {
			return err
		}

		return callback(result)
	})

	if err != nil {
		return Member{}, err
	}

	return result, nil
}

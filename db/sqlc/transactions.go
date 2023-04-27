package db

import (
	"context"
	"time"

	"github.com/google/uuid"
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

type CreateMemberTxParams struct {
	Username  string
	Fullname  string
	Email     string
	Password  string
	Plan      int32
	CreatedAt time.Time
	ExpiredAt time.Time
	AutoRenew bool
}

type CreateSessionTxParams struct {
	ID        uuid.UUID
	Token     string
	ClientIp  string
	UserAgent string
	IsActive  bool
	CreatedAt time.Time
	ExpiredAt time.Time
}

type CreateMemberResult struct {
	session Session
	member  Member
}

func (store *SQLStore) CreateMemberTx(ctx context.Context, memberArg CreateMemberTxParams, SessionArg CreateSessionTxParams) (*CreateMemberResult, error) {
	var result CreateMemberResult

	err := store.execTx(ctx, func(q *Queries) error {

		member, err := q.CreateMember(ctx, CreateMemberParams(memberArg))
		if err != nil {
			return err
		}

		session, err := q.CreateSession(ctx, CreateSessionParams{
			ID:        SessionArg.ID,
			MemberID:  member.ID,
			Token:     SessionArg.Token,
			ClientIp:  SessionArg.ClientIp,
			UserAgent: SessionArg.UserAgent,
			IsActive:  SessionArg.IsActive,
			CreatedAt: SessionArg.CreatedAt,
			ExpiredAt: SessionArg.ExpiredAt,
		})
		if err != nil {
			return err
		}

		result.member = member
		result.session = session

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

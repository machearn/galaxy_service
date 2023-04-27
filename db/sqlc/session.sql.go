// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO "sessions" (
  "id", "member_id", "token", "client_ip", "user_agent", "is_active", "created_at", "expired_at"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, member_id, token, client_ip, user_agent, is_active, created_at, expired_at
`

type CreateSessionParams struct {
	ID        uuid.UUID `json:"id"`
	MemberID  int32     `json:"member_id"`
	Token     string    `json:"token"`
	ClientIp  string    `json:"client_ip"`
	UserAgent string    `json:"user_agent"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.MemberID,
		arg.Token,
		arg.ClientIp,
		arg.UserAgent,
		arg.IsActive,
		arg.CreatedAt,
		arg.ExpiredAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.MemberID,
		&i.Token,
		&i.ClientIp,
		&i.UserAgent,
		&i.IsActive,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}
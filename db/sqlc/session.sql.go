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
  "id", "member_id", "refresh_token", "client_ip", "user_agent", "created_at", "expired_at"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING id, member_id, refresh_token, client_ip, user_agent, is_blocked, created_at, expired_at
`

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	MemberID     int32     `json:"member_id"`
	RefreshToken string    `json:"refresh_token"`
	ClientIp     string    `json:"client_ip"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.MemberID,
		arg.RefreshToken,
		arg.ClientIp,
		arg.UserAgent,
		arg.CreatedAt,
		arg.ExpiredAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.MemberID,
		&i.RefreshToken,
		&i.ClientIp,
		&i.UserAgent,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id, member_id, refresh_token, client_ip, user_agent, is_blocked, created_at, expired_at FROM "sessions"
WHERE "id" = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.MemberID,
		&i.RefreshToken,
		&i.ClientIp,
		&i.UserAgent,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

-- name: CreateSession :one
INSERT INTO "sessions" (
  "id", "member_id", "token", "client_ip", "user_agent", "is_active", "created_at", "expired_at"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;
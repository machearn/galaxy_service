-- name: CreateSession :one
INSERT INTO "sessions" (
  "id", "member_id", "refresh_token", "client_ip", "user_agent", "created_at", "expired_at"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM "sessions"
WHERE "id" = $1 LIMIT 1;
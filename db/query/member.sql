-- name: GetMember :one
SELECT * FROM members WHERE id = $1 LIMIT 1;

-- name: GetMemberByName :one
SELECT * FROM members WHERE username = $1 LIMIT 1;

-- name: ListMembers :many
SELECT * FROM members
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateMember :one
INSERT INTO members (
  username, fullname, email, password, plan, created_at, expired_at, auto_renew
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateMember :one
UPDATE members SET
  username = coalesce(sqlc.narg('username'), username),
  fullname = coalesce(sqlc.narg('fullname'), fullname),
  email = coalesce(sqlc.narg('email'), email),
  password = coalesce(sqlc.narg('password'), password),
  plan = coalesce(sqlc.narg('plan'), plan),
  expired_at = coalesce(sqlc.narg('expired_at'), expired_at),
  auto_renew = coalesce(sqlc.narg('auto_renew'), auto_renew)
WHERE id = @id
RETURNING *;

-- name: DeleteMember :exec
DELETE FROM members WHERE id = $1;
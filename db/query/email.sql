-- name: CreateVerificationEmail :one
INSERT INTO verification_emails(
  member_id, email, secret_code
) VALUES (
  $1, $2, $3
) RETURNING *;
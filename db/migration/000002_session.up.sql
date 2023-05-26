ALTER TABLE "members" ADD COLUMN "password" varchar NOT NULL;

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "member_id" integer NOT NULL,
  "refresh_token" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL,
  "expired_at" timestamp NOT NULL
);
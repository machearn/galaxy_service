ALTER TABLE "members" ADD COLUMN "password" varchar NOT NULL;

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "member_id" integer NOT NULL,
  "token" varchar UNIQUE NOT NULL,
  "client_ip" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "is_active" boolean NOT NULL,
  "created_at" timestamp NOT NULL,
  "expired_at" timestamp NOT NULL
);
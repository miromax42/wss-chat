-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2022-04-19T08:45:10.563Z

CREATE TABLE "rooms" (
  "name" varchar PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "messages" (
  "id" SERIAL PRIMARY KEY,
  "sender" varchar NOT NULL,
  "room" varchar NOT NULL,
  "payload" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "messages" ("created_at", "room");

ALTER TABLE "messages" ADD FOREIGN KEY ("room") REFERENCES "rooms" ("name");

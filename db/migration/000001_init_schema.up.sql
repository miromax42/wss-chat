CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "rooms" (
  "name" varchar PRIMARY KEY,
  "creator" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "messages" (
  "id" SERIAL PRIMARY KEY,
  "sender" varchar NOT NULL,
  "room" varchar NOT NULL,
  "payload" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "rooms" ADD FOREIGN KEY ("creator") REFERENCES "users" ("username");

ALTER TABLE "messages" ADD FOREIGN KEY ("sender") REFERENCES "users" ("username");

ALTER TABLE "messages" ADD FOREIGN KEY ("room") REFERENCES "rooms" ("name");

CREATE INDEX ON "messages" ("created_at", "room");

CREATE TABLE "users" (
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "createdAt" timestamptz NOT NULL DEFAULT now(),
    "updatedAt" timestamptz,
    "ID" varchar PRIMARY KEY NOT NULL,
    "admin" boolean NOT NULL DEFAULT FALSE
);

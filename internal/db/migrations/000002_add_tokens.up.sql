CREATE TABLE "tokens" (
                          "userID" varchar NOT NULL,
                          "type" varchar NOT NULL DEFAULT 'login',
                          "ID" varchar PRIMARY KEY NOT NULL,
                          "createdAt" timestamptz
);

ALTER TABLE "tokens" ADD CONSTRAINT token_user_id FOREIGN KEY ("userID") REFERENCES "users" ("ID");

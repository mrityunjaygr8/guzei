CREATE TABLE "hatinIt" (
                           "imageURI" varchar NOT NULL,
                           "title" varchar NOT NULL,
                           "userID" varchar NOT NULL,
                           "createdAt" timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE "hatinIt" ADD CONSTRAINT "hatinIt_user_id" FOREIGN KEY ("userID") REFERENCES "users" ("ID");
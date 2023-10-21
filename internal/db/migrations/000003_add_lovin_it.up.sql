CREATE TABLE "lovinIt" (
                           "byline" varchar NOT NULL,
                           "detailsURI" varchar NOT NULL,
                           "imageURI" varchar NOT NULL,
                           "thumbURI" varchar NOT NULL,
                           "title" varchar NOT NULL,
                           "userID" varchar NOT NULL,
                           "lastUsed" date NOT NULL DEFAULT CURRENT_DATE,
                           "createdAt" timestamptz
);

ALTER TABLE "lovinIt" ADD CONSTRAINT "lovinIt_user_id" FOREIGN KEY ("userID") REFERENCES "users" ("ID");
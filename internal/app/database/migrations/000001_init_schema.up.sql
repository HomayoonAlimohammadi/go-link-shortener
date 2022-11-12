CREATE TABLE "linkshortener" (
    "url" VARCHAR(512) UNIQUE NOT NULL,
    "token" VARCHAR(10) UNIQUE NOT NULL
);

CREATE INDEX ON "linkshortener" ("url");

CREATE INDEX ON "linkshortener" ("token");
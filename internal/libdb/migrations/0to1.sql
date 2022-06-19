-- Create new Surchable database from scratch
-- v0 -> v1

CREATE TABLE "domain_queue"
(
    id         UUID                NOT NULL PRIMARY KEY,
    created_at TIMESTAMP           NOT NULL DEFAULT now(),
    domain     VARCHAR(253) UNIQUE NOT NULL
);

CREATE TABLE "page_loads"
(
    id              UUID        NOT NULL PRIMARY KEY,
    url             TEXT UNIQUE NOT NULL,
    normalised_url  TEXT UNIQUE NOT NULL
    content_sha1    VARCHAR(40),
    loaded_at       TIMESTAMP   NOT NULL DEFAULT now(),
    not_load_before TIMESTAMP
);

CREATE TABLE "current_jobs"
(
    id            UUID        NOT NULL PRIMARY KEY,
    queue_item    UUID        NOT NULL REFERENCES "domain_queue" (id),
    worker_id     TEXT UNIQUE NOT NULL,
    last_check_in TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE TABLE "version"
(
    version INT
);

INSERT INTO "version"(version)
VALUES (1);
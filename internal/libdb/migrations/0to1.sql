-- Create new Surchable database from scratch
-- v0 -> v1

CREATE TABLE "domain_queue"
(
    id         uuid      not null primary key,
    created_at timestamp not null default now(),
    domain     varchar(253),
    subdomain  varchar(253)
);

CREATE TABLE "page_loads"
(
    id              UUID      NOT NULL PRIMARY KEY,
    url             TEXT      NOT NULL,
    content         TEXT,
    loaded_at       TIMESTAMP NOT NULL DEFAULT now(),
    not_load_before TIMESTAMP
);

CREATE TABLE "current_jobs"
(
    id            UUID      NOT NULL PRIMARY KEY,
    queue_item    UUID      NOT NULL REFERENCES domain_queue (id),
    worker_id     UUID      NOT NULL,
    last_check_in TIMESTAMP NOT NULL
);

CREATE TABLE "version"
(
    version INT
);

INSERT INTO "version"(version)
VALUES (1);
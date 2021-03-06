-- Create new Surchable database from scratch
-- v0 -> v1

CREATE TABLE "domain_queue"
(
    id         UUID                NOT NULL PRIMARY KEY,
    created_at TIMESTAMP           NOT NULL DEFAULT now(),
    domain     VARCHAR(253) UNIQUE NOT NULL,
    start      TEXT                NOT NULL DEFAULT '/',
    priority   INT                 NOT NULL DEFAULT 0
);

CREATE TABLE "domain_blocklist"
(
    domain TEXT NOT NULL PRIMARY KEY,
    reason TEXT NOT NULL
);

CREATE TABLE "current_jobs"
(
    id            UUID        NOT NULL PRIMARY KEY,
    queue_item    UUID        NOT NULL REFERENCES "domain_queue" (id),
    crawler_id    TEXT UNIQUE NOT NULL,
    last_check_in TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE TABLE "page_loads"
(
    id              UUID        NOT NULL PRIMARY KEY,
    url             TEXT UNIQUE NOT NULL,
    normalised_url  TEXT UNIQUE NOT NULL,
    loaded_at       TIMESTAMP   NOT NULL DEFAULT now(),
    not_load_before TIMESTAMP
);

CREATE TABLE "page_information"
(
    id                         UUID        NOT NULL PRIMARY KEY,
    load_id                    UUID UNIQUE NOT NULL REFERENCES "page_loads" (id),
    page_title                 TEXT,
    page_meta_description_text TEXT,
    page_content_text          TEXT,
    page_raw_html              TEXT        NOT NULL,
    raw_html_sha1              VARCHAR(40) NOT NULL,
    outbound_links             TEXT ARRAY
);

CREATE TABLE "search_index"
(
    token          TEXT   NOT NULL,
    page_id        UUID   NOT NULL REFERENCES "page_information" (id),
    classification BIT(5) NOT NULL
);

CREATE TABLE "token_frequencies"
(
    token          TEXT   NOT NULL,
    page_id        UUID   NOT NULL REFERENCES "page_information" (id),
    frequency      INT    NOT NULL,
    classification BIT(5) NOT NULL
);

CREATE TABLE "version"
(
    version INT
);

INSERT INTO "version"(version)
VALUES (1);
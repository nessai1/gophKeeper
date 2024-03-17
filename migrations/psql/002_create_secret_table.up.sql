BEGIN;

CREATE TABLE IF NOT EXISTS secret_metadata (
    uuid uuid not null primary key,
    owner_uuid uuid not null,
    name varchar(255) not null,
    type smallint not null,
    created timestamp not null default now(),
    updated timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS plain_secret (
    uuid uuid primary key,
    data text not null,
    FOREIGN KEY (uuid) REFERENCES plain_secret (uuid)
);

CREATE INDEX secret_metadata_owner_uuid ON secret_metadata (owner_uuid);
CREATE INDEX secret_metadata_owner_uuid_name ON secret_metadata (owner_uuid, name);
COMMIT;

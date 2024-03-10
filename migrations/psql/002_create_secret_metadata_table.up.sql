BEGIN;
CREATE TABLE IF NOT EXISTS secret_metadata (
    uuid uuid not null primary key,
    owner_uuid uuid not null,
    name varchar(255) not null,
    type smallint not null
);
CREATE INDEX secret_metadata_owner_uuid ON secret_metadata (owner_uuid);
COMMIT;

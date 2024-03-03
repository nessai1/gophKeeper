BEGIN;
CREATE TABLE IF NOT EXISTS storage_objects (
    uuid uuid not null primary key,
    owner_id int not null,
    original_name varchar(255) not null,
    type varchar(50) not null
);
CREATE INDEX storage_objects_owner_id ON storage_objects (owner_id);
COMMIT;

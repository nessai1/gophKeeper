BEGIN;
DROP INDEX IF EXISTS secret_metadata_owner_uuid;
DROP INDEX IF EXISTS secret_metadata_owner_uuid_name;
DROP TABLE IF EXISTS plain_secret;
DROP TABLE IF EXISTS secret_metadata;
COMMIT;

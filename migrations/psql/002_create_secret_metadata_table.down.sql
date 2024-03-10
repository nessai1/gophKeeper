BEGIN;
DROP INDEX IF EXISTS secret_metadata_owner_uuid;
DROP TABLE IF EXISTS secret_metadata;
COMMIT;

DROP INDEX IF EXISTS idx_user_id_type_name;
DROP INDEX IF EXISTS idx_secrets_created_at;
DROP INDEX IF EXISTS idx_secrets_versions_secret_id;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS secrets;
DROP TABLE IF EXISTS secret_versions;

DROP TYPE IF EXISTS secret_type;

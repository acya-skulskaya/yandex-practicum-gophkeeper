CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login         VARCHAR(255) UNIQUE NOT NULL,
    password      VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TYPE IF EXISTS secret_type;
CREATE TYPE secret_type AS ENUM ('login_password', 'text', 'binary', 'bank_card');

CREATE TABLE IF NOT EXISTS secrets (
    id            INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id           INTEGER        NOT NULL,
    type            secret_type,
    name      VARCHAR(255)  NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS secret_versions (
    id            INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    secret_id           INTEGER        NOT NULL,
    data TEXT NULL ,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_user_id_type_name ON secrets(user_id, type, name);
CREATE INDEX idx_secrets_created_at ON secrets (created_at);
CREATE INDEX idx_secrets_versions_secret_id ON secret_versions (secret_id);

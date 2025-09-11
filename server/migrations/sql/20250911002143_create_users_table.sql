-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    user_id                     VARCHAR(100) PRIMARY KEY      DEFAULT uuid_generate_v4(),
    name                        VARCHAR(100) NOT NULL DEFAULT '',
    batch                       INT          DEFAULT 0,
    email                       VARCHAR(255) NOT NULL DEFAULT '',
    password                    VARCHAR(255) NOT NULL DEFAULT '',
    role                        VARCHAR(255) NOT NULL DEFAULT '',
    created_at                  TIMESTAMP    DEFAULT current_timestamp,
    created_by                  VARCHAR(100) NOT NULL DEFAULT '',
    updated_at                  TIMESTAMP    DEFAULT current_timestamp,
    updated_by                  VARCHAR(100) NOT NULL DEFAULT '',
    deleted_at                  TIMESTAMP    DEFAULT NULL,
    deleted_by                  VARCHAR(100) NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

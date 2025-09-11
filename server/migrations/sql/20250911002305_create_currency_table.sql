-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currency
(
    id                          SERIAL PRIMARY KEY,
    name                        VARCHAR(100) NOT NULL DEFAULT '',
    exchange_rate               INT          NOT NULL DEFAULT 0,
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
DROP TABLE IF EXISTS currency;
-- +goose StatementEnd

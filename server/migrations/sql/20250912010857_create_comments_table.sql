-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments
(
    comment_id                  VARCHAR(100) PRIMARY KEY      DEFAULT uuid_generate_v4(),
    serial_number               SERIAL,
    name                        VARCHAR(100) NOT NULL DEFAULT '',
    email                       VARCHAR(100) NOT NULL DEFAULT '',
    comment                     TEXT NOT NULL DEFAULT '',
    created_at                  TIMESTAMP    DEFAULT current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd

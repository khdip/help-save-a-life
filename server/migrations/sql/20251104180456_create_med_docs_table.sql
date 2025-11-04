-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS med_docs
(
    id                          VARCHAR(100) PRIMARY KEY      DEFAULT uuid_generate_v4(),
    serial_number               SERIAL,
    name                        VARCHAR(100) NOT NULL DEFAULT '',
    type                        VARCHAR(10)  NOT NULL DEFAULT '',
    uploaded_at                 TIMESTAMP    DEFAULT current_timestamp,
    uploaded_by                 VARCHAR(100) NOT NULL DEFAULT '',
    deleted_at                  TIMESTAMP    DEFAULT NULL,
    deleted_by                  VARCHAR(100) NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS med_docs;
-- +goose StatementEnd

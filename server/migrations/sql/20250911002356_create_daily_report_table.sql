-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS daily_report
(
    report_id                   SERIAL PRIMARY KEY,
    date                        VARCHAR(100) NOT NULL DEFAULT '', 
    amount                      INT NOT NULL DEFAULT 0,
    currency                    VARCHAR(10) NOT NULL DEFAULT 'BDT',
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
DROP TABLE IF EXISTS daily_report;
-- +goose StatementEnd

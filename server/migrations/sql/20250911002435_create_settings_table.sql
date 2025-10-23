-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS settings
(
    patient_name                    VARCHAR(20) NOT NULL DEFAULT '', 
    title                           VARCHAR(30) NOT NULL DEFAULT '', 
    banner_title                    VARCHAR(100) NOT NULL DEFAULT '',
    highlighted_banner_title        VARCHAR(100) NOT NULL DEFAULT '',
    banner_description              TEXT NOT NULL DEFAULT '',
    highlighted_banner_description  TEXT NOT NULL DEFAULT '',
    banner_image                    VARCHAR(100) NOT NULL DEFAULT '',
    about_patient                   TEXT NOT NULL DEFAULT '',
    target_amount                   INT NOT NULL DEFAULT 0,
    show_med_docs                   BOOLEAN NOT NULL DEFAULT true,
    show_collection                 BOOLEAN NOT NULL DEFAULT true,
    show_daily_report               BOOLEAN NOT NULL DEFAULT true,
    show_fund_updates               BOOLEAN NOT NULL DEFAULT true,
    calculate_collection            INT NOT NULL DEFAULT 0,
    total_amount                    INT NOT NUll DEFAULT 0,
    updated_at                      TIMESTAMP    DEFAULT current_timestamp,
    updated_by                      VARCHAR(100) NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS settings;
-- +goose StatementEnd

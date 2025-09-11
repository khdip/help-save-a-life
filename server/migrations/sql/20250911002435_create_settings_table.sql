-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS settings
(
    patient_name                VARCHAR(100) NOT NULL DEFAULT '', 
    target_amount               INT NOT NULL DEFAULT 0,
    show_med_docs               BOOLEAN NOT NULL DEFAULT true,
    show_collection             BOOLEAN NOT NULL DEFAULT true,
    show_daily_report           BOOLEAN NOT NULL DEFAULT true,
    show_fund_updates           BOOLEAN NOT NULL DEFAULT true
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS settings;
-- +goose StatementEnd

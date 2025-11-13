package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/khdip/help-save-a-life/server/storage"
)

const getSettings = `SELECT * FROM settings;`

func (s *Storage) GetSettings(ctx context.Context, sett storage.Settings) (*storage.Settings, error) {
	var res storage.Settings
	if err := s.db.Get(&res, getSettings); err != nil {
		return nil, fmt.Errorf("executing settings details: %w", err)
	}
	return &res, nil
}

const updateSettings = `
UPDATE
	 settings
SET
	patient_name = :patient_name,
    title = :title,
	banner_title = :banner_title,
	highlighted_banner_title = :highlighted_banner_title,
	banner_description = :banner_description,
	highlighted_banner_description = :highlighted_banner_description,
	banner_image = :banner_image,
	about_patient = :about_patient,
	target_amount = :target_amount,
	show_med_docs = :show_med_docs,
	show_collection = :show_collection,
	show_daily_report = :show_daily_report,
	show_fund_updates = :show_fund_updates,
	calculate_collection = :calculate_collection,
	total_amount = :total_amount,
	updated_at = now(),
	updated_by = :updated_by
RETURNING 
	updated_at;
`

func (s *Storage) UpdateSettings(ctx context.Context, sett storage.Settings) (*storage.Settings, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateSettings)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&sett, sett); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing settings update: %w", err)
	}
	return &sett, nil
}

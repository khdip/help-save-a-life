package storage

import (
	"database/sql"
	"time"
)

type User struct {
	UserID       string `db:"user_id"`
	SerialNumber int32  `db:"serial_number"`
	Name         string `db:"name"`
	Batch        int32  `db:"batch"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	Role         string `db:"role"`
	CRUDTimeDate
}

type Collection struct {
	CollectionID  string `db:"collection_id"`
	SerialNumber  int32  `db:"serial_number"`
	AccountType   string `db:"account_type"`
	AccountNumber string `db:"account_number"`
	Sender        string `db:"sender"`
	Date          string `db:"date"`
	Amount        int32  `db:"amount"`
	Currency      string `db:"currency"`
	CRUDTimeDate
}

type DailyReport struct {
	ReportID     string `db:"report_id"`
	SerialNumber int32  `db:"serial_number"`
	Date         string `db:"date"`
	Amount       int32  `db:"amount"`
	Currency     string `db:"currency"`
	CRUDTimeDate
}

type Comment struct {
	CommentID    string    `db:"comment_id"`
	SerialNumber int32     `db:"serial_number"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Comment      string    `db:"comment"`
	CreatedAt    time.Time `db:"created_at,omitempty"`
}

type Currency struct {
	ID           string `db:"id"`
	SerialNumber int32  `db:"serial_number"`
	Name         string `db:"name"`
	ExchangeRate int32  `db:"exchange_rate"`
	CRUDTimeDate
}

type Settings struct {
	PatientName          string `db:"patient_name"`
	TargetAmount         int32  `db:"target_amount"`
	ShowMedicalDocuments bool   `db:"show_med_docs"`
	ShowCollection       bool   `db:"show_collection"`
	ShowDailyReport      bool   `db:"show_daily_report"`
	ShowFundUpdates      bool   `db:"show_fund_updates"`
	CRUDTimeDate
}

type CRUDTimeDate struct {
	CreatedAt time.Time      `db:"created_at,omitempty"`
	CreatedBy string         `db:"created_by"`
	UpdatedAt time.Time      `db:"updated_at,omitempty"`
	UpdatedBy string         `db:"updated_by,omitempty"`
	DeletedAt sql.NullTime   `db:"deleted_at,omitempty"`
	DeletedBy sql.NullString `db:"deleted_by,omitempty"`
}

type Filter struct {
	Offset     int32
	Limit      int32
	SortBy     string
	Order      string
	SearchTerm string
}

type Stats struct {
	Count       int32
	TotalAmount int32 `db:"coalesce"`
}

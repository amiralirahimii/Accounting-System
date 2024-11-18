package models

import (
	"database/sql"
	"time"
)

type VoucherItem struct {
	ID         int           `gorm:"primaryKey"`
	VoucherID  int           `gorm:"not null"`
	SLID       int           `gorm:"not null"`
	DLID       sql.NullInt64 `gorm:""`
	Debit      int           `gorm:"check:debit >= 0"`
	Credit     int           `gorm:"check:credit >= 0"`
	RowVersion int           `gorm:"not null;default:0"`
	CreatedAt  time.Time     `gorm:"autoCreateTime"`
	UpdatedAt  time.Time     `gorm:"autoUpdateTime"`
}

func (VoucherItem) TableName() string {
	return "voucher_item"
}

package models

import (
	"database/sql"
	"time"
)

type VoucherItem struct {
	ID        int
	VoucherID int
	SLID      int
	DLID      sql.NullInt64
	Debit     int
	Credit    int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (VoucherItem) TableName() string {
	return "voucher_item"
}

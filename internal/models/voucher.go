package models

import "time"

type Voucher struct {
	ID         int
	Number     string
	RowVersion int
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (Voucher) TableName() string {
	return "voucher"
}

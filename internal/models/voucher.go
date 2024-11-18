package models

import "time"

type Voucher struct {
	ID         int       `gorm:"primaryKey"`
	Number     string    `gorm:"size:64;not null;unique"`
	RowVersion int       `gorm:"not null;default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (Voucher) TableName() string {
	return "voucher"
}

package models

import "time"

type DL struct {
	ID         int       `gorm:"primaryKey"`
	Code       string    `gorm:"size:64;not null;unique"`
	Title      string    `gorm:"size:64;not null;unique"`
	RowVersion int       `gorm:"not null;default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (DL) TableName() string {
	return "dl"
}

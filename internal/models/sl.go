package models

import "time"

type SL struct {
	ID         int       `gorm:"primaryKey"`
	Code       string    `gorm:"size:64;not null;unique"`
	Title      string    `gorm:"size:64;not null;unique"`
	IsDL       bool      `gorm:"not null"`
	RowVersion int       `gorm:"not null;default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (SL) TableName() string {
	return "sl"
}

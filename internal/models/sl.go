package models

import "time"

type SL struct {
	ID         int
	Code       string
	Title      string
	HasDL      bool
	RowVersion int
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (SL) TableName() string {
	return "sl"
}

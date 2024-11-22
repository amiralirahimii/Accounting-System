package models

import "time"

type DL struct {
	ID         int
	Code       string
	Title      string
	RowVersion int
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (DL) TableName() string {
	return "dl"
}

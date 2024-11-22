package mappers

import (
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/models"
)

func ToSlDto(sl *models.SL) *dtos.SLDto {
	return &dtos.SLDto{
		ID:         sl.ID,
		Code:       sl.Code,
		Title:      sl.Title,
		RowVersion: sl.RowVersion,
	}
}

package mappers

import (
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/models"
)

func ToDLDto(dl *models.DL) *dtos.DLDto {
	return &dtos.DLDto{
		ID:         dl.ID,
		Code:       dl.Code,
		Title:      dl.Title,
		RowVersion: dl.RowVersion,
	}
}

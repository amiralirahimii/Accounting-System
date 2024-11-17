package services

import (
	"accountingsystem/db"
	"accountingsystem/internal/constants"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/sl"
	"errors"
	"log"

	"gorm.io/gorm"
)

type SLService struct{}

func (s *SLService) CreateSL(req *sl.InsertRequest) (*models.SL, error) {
	if req.Code == "" || len(req.Code) > 64 {
		return nil, constants.ErrCodeEmptyOrTooLong
	}
	if req.Title == "" || len(req.Title) > 64 {
		return nil, constants.ErrTitleEmptyOrTooLong
	}

	var existingSL models.SL
	if err := db.DB.Where("code = ? OR title = ?", req.Code, req.Title).First(&existingSL).Error; err == nil {
		if existingSL.Code == req.Code {
			return nil, constants.ErrCodeAlreadyExists
		}
		if existingSL.Title == req.Title {
			return nil, constants.ErrTitleAlreadyExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Unexpected error while checking for duplicates: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	sl := models.SL{
		Code:       req.Code,
		Title:      req.Title,
		HasDL:      req.HasDL,
		RowVersion: 0,
	}

	if err := db.DB.Create(&sl).Error; err != nil {
		return nil, err
	}

	return &sl, nil
}

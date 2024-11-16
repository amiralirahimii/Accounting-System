package services

import (
	"accountingsystem/db"
	"accountingsystem/internal/constants"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"errors"
	"log"

	"gorm.io/gorm"
)

type DLService struct{}

func (s *DLService) CreateDL(req *dl.InsertRequest) (*models.DL, error) {
	if req.Code == "" || len(req.Code) > 64 {
		return nil, constants.ErrCodeEmptyOrTooLong
	}
	if req.Title == "" || len(req.Title) > 64 {
		return nil, constants.ErrTitleEmptyOrTooLong
	}

	var existingDL models.DL
	if err := db.DB.Where("code = ? OR title = ?", req.Code, req.Title).First(&existingDL).Error; err == nil {
		if existingDL.Code == req.Code {
			return nil, constants.ErrCodeAlreadyExists
		}
		if existingDL.Title == req.Title {
			return nil, constants.ErrTitleAlreadyExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("unexpected error while checking for duplicates: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	dl := models.DL{
		Code:       req.Code,
		Title:      req.Title,
		RowVersion: 0,
	}

	if err := db.DB.Create(&dl).Error; err != nil {
		return nil, err
	}

	return &dl, nil
}

func (s *DLService) UpdateDL(req *dl.UpdateRequest) (*models.DL, error) {
	if req.Code == "" || len(req.Code) > 64 {
		return nil, constants.ErrCodeEmptyOrTooLong
	}
	if req.Title == "" || len(req.Title) > 64 {
		return nil, constants.ErrTitleEmptyOrTooLong
	}

	var targerDL models.DL
	if err := db.DB.Find(&targerDL, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrDLNotFound
		}
		log.Printf("unexpected error while checking for duplicates: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if targerDL.RowVersion != req.Version {
		return nil, constants.ErrVersionOutdated
	}

	var existingDL models.DL
	if err := db.DB.Where("(code = ? OR title = ?) AND id != ?", req.Code, req.Title, req.ID).First(&existingDL).Error; err == nil {
		if existingDL.Code == req.Code {
			return nil, constants.ErrCodeAlreadyExists
		}
		if existingDL.Title == req.Title {
			return nil, constants.ErrTitleAlreadyExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("unexpected error while checking for duplicates: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	targerDL.Code = req.Code
	targerDL.Title = req.Title
	targerDL.RowVersion++
	if err := db.DB.Save(&targerDL).Error; err != nil {
		return nil, err
	}

	return &targerDL, nil
}

package services

import (
	"accountingsystem/db"
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"errors"
	"log"

	"gorm.io/gorm"
)

type DLService struct{}

func (s *DLService) CreateDL(req *dl.InsertRequest) (*dtos.DLDto, error) {
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

	return mappers.ToDLDto(&dl), nil
}

func (s *DLService) UpdateDL(req *dl.UpdateRequest) (*dtos.DLDto, error) {
	if req.Code == "" || len(req.Code) > 64 {
		return nil, constants.ErrCodeEmptyOrTooLong
	}
	if req.Title == "" || len(req.Title) > 64 {
		return nil, constants.ErrTitleEmptyOrTooLong
	}

	var targetDL models.DL
	if err := db.DB.Where("id = ?", req.ID).First(&targetDL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrDLNotFound
		}
		log.Printf("unexpected error while finding DL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if targetDL.RowVersion != req.Version {
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

	targetDL.Code = req.Code
	targetDL.Title = req.Title
	targetDL.RowVersion++
	if err := db.DB.Save(&targetDL).Error; err != nil {
		return nil, err
	}

	return mappers.ToDLDto(&targetDL), nil
}

func (s *DLService) DeleteDL(req *dl.DeleteRequest) error {
	var targetDL models.DL
	if err := db.DB.Where("id = ?", req.ID).First(&targetDL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrDLNotFound
		}
		log.Printf("unexpected error while finding DL: %v", err)
		return constants.ErrUnexpectedError
	}

	if targetDL.RowVersion != req.Version {
		return constants.ErrVersionOutdated
	}

	var VoucherItemRefrencingThisDL models.VoucherItem
	if err := db.DB.Where("dl_id = ?", targetDL.ID).First(&VoucherItemRefrencingThisDL).Error; err == nil {
		return constants.ErrThereIsRefrenceToDL
	}

	if err := db.DB.Delete(&targetDL).Error; err != nil {
		log.Printf("unexpected error while deleting DL: %v", err)
		return constants.ErrUnexpectedError
	}

	return nil
}

func (s *DLService) GetDL(req *dl.GetRequest) (*dtos.DLDto, error) {
	var targetDL models.DL
	if err := db.DB.Where("id = ?", req.ID).First(&targetDL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrDLNotFound
		}
		log.Printf("unexpected error while finding DL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return mappers.ToDLDto(&targetDL), nil
}

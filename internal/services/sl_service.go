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

func (s *SLService) UpdateSL(req *sl.UpdateRequest) (*models.SL, error) {
	if req.Code == "" || len(req.Code) > 64 {
		return nil, constants.ErrCodeEmptyOrTooLong
	}
	if req.Title == "" || len(req.Title) > 64 {
		return nil, constants.ErrTitleEmptyOrTooLong
	}

	var targetSL models.SL
	if err := db.DB.Where("id = ?", req.ID).First(&targetSL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrSLNotFound
		}
		log.Printf("unexpected error while finding DL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if targetSL.RowVersion != req.Version {
		return nil, constants.ErrVersionOutdated
	}

	var VoucherItemRefrencingThisSL models.VoucherItem
	if err := db.DB.Where("sl_id = ?", targetSL.ID).First(&VoucherItemRefrencingThisSL).Error; err == nil {
		return nil, constants.ErrThereIsRefrenceToSL
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

	targetSL.Code = req.Code
	targetSL.Title = req.Title
	targetSL.HasDL = req.HasDL
	targetSL.RowVersion++

	if err := db.DB.Save(&targetSL).Error; err != nil {
		return nil, err
	}

	return &targetSL, nil
}

// TODO check for any refrences here before deleting
func (s *SLService) DeleteSL(req *sl.DeleteRequest) error {
	var targetSL models.SL
	if err := db.DB.Where("id = ?", req.ID).First(&targetSL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrSLNotFound
		}
		log.Printf("unexpected error while finding SL: %v", err)
		return constants.ErrUnexpectedError
	}

	if targetSL.RowVersion != req.Version {
		return constants.ErrVersionOutdated
	}

	if err := db.DB.Delete(&targetSL).Error; err != nil {
		return err
	}

	return nil
}

func (s *SLService) GetSL(req *sl.GetRequest) (*models.SL, error) {
	var sl models.SL
	if err := db.DB.Where("id = ?", req.ID).First(&sl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrSLNotFound
		}
		log.Printf("unexpected error while finding SL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return &sl, nil
}

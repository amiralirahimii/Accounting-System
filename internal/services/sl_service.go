package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/sl"
	"errors"
	"log"

	"gorm.io/gorm"
)

type SLService struct {
	db *gorm.DB
}

func (s *SLService) InitService(db *gorm.DB) {
	s.db = db
}

func (s *SLService) CreateSL(req *sl.InsertRequest) (*dtos.SLDto, error) {
	if err := s.validateSLInsertRequest(req); err != nil {
		return nil, err
	}

	sl := models.SL{
		Code:       req.Code,
		Title:      req.Title,
		HasDL:      req.HasDL,
		RowVersion: 0,
	}

	if err := s.db.Create(&sl).Error; err != nil {
		return nil, err
	}

	return mappers.ToSlDto(&sl), nil
}

func (s *SLService) validateCodeAndTitleUnique(code string, title string) error {
	var existingSL models.SL
	if err := s.db.Where("code = ? OR title = ?", code, title).First(&existingSL).Error; err == nil {
		if existingSL.Code == code {
			return constants.ErrCodeAlreadyExists
		}
		if existingSL.Title == title {
			return constants.ErrTitleAlreadyExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("unexpected error while checking for duplicates: %v", err)
		return constants.ErrUnexpectedError
	}
	return nil
}

func (s *SLService) validateSLInsertRequest(req *sl.InsertRequest) error {
	if err := s.validateCodeAndTitleLength(req.Code, req.Title); err != nil {
		return err
	}
	if err := s.validateCodeAndTitleUnique(req.Code, req.Title); err != nil {
		return err
	}
	return nil
}

func (s *SLService) validateCodeAndTitleLength(code string, title string) error {
	if code == "" || len(code) > 64 {
		return constants.ErrCodeEmptyOrTooLong
	}
	if title == "" || len(title) > 64 {
		return constants.ErrTitleEmptyOrTooLong
	}
	return nil
}

func (s *SLService) UpdateSL(req *sl.UpdateRequest) (*dtos.SLDto, error) {
	targetSL, err := s.validateSLUpdateRequest(req)
	if err != nil {
		return nil, err
	}

	targetSL.Code = req.Code
	targetSL.Title = req.Title
	targetSL.HasDL = req.HasDL
	targetSL.RowVersion++

	if err := s.db.Save(targetSL).Error; err != nil {
		return nil, err
	}

	return mappers.ToSlDto(targetSL), nil
}

func (s *SLService) validateSLUpdateRequest(req *sl.UpdateRequest) (*models.SL, error) {
	if err := s.validateCodeAndTitleLength(req.Code, req.Title); err != nil {
		return nil, err
	}
	targetSL, err := s.validateSLExists(req.ID)
	if err != nil {
		return nil, err
	}
	if err := s.validateVersion(req.Version, targetSL.RowVersion); err != nil {
		return nil, err
	}
	if err := s.validateSLHasNoReferences(req.ID); err != nil {
		return nil, err
	}
	if err := s.validateCodeAndTitleUniqueWithDifferentId(req.Code, req.Title, req.ID); err != nil {
		return nil, err
	}
	return targetSL, nil
}

func (s *SLService) validateCodeAndTitleUniqueWithDifferentId(code string, title string, id int) error {
	var existingSL models.SL
	if err := s.db.Where("(code = ? OR title = ?) AND id != ?", code, title, id).First(&existingSL).Error; err == nil {
		if existingSL.Code == code {
			return constants.ErrCodeAlreadyExists
		}
		if existingSL.Title == title {
			return constants.ErrTitleAlreadyExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("unexpected error while checking for duplicates: %v", err)
		return constants.ErrUnexpectedError
	}
	return nil
}

func (s *SLService) validateSLHasNoReferences(id int) error {
	var VoucherItemRefrencingThisSL models.VoucherItem
	if err := s.db.Where("sl_id = ?", id).First(&VoucherItemRefrencingThisSL).Error; err == nil {
		return constants.ErrThereIsRefrenceToSL
	}
	return nil
}

func (s *SLService) validateVersion(reqVersion int, targetVersion int) error {
	if reqVersion != targetVersion {
		return constants.ErrVersionOutdated
	}
	return nil
}

func (s *SLService) validateSLExists(id int) (*models.SL, error) {
	var targetSL models.SL
	if err := s.db.Where("id = ?", id).First(&targetSL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrSLNotFound
		}
		log.Printf("unexpected error while finding SL: %v", err)
		return nil, constants.ErrUnexpectedError
	}
	return &targetSL, nil
}

func (s *SLService) DeleteSL(req *sl.DeleteRequest) error {
	targetSL, err := s.validateSLDeleteRequest(req)
	if err != nil {
		return err
	}

	if err := s.db.Delete(&targetSL).Error; err != nil {
		return err
	}

	return nil
}

func (s *SLService) validateSLDeleteRequest(req *sl.DeleteRequest) (*models.SL, error) {
	targetSL, err := s.validateSLExists(req.ID)
	if err != nil {
		return nil, err
	}
	if err := s.validateVersion(req.Version, targetSL.RowVersion); err != nil {
		return nil, err
	}
	if err := s.validateSLHasNoReferences(req.ID); err != nil {
		return nil, err
	}
	return targetSL, nil
}

func (s *SLService) GetSL(req *sl.GetRequest) (*dtos.SLDto, error) {
	targetSL, err := s.validateSLGetRequest(req)
	if err != nil {
		return nil, err
	}

	return mappers.ToSlDto(targetSL), nil
}

func (s *SLService) validateSLGetRequest(req *sl.GetRequest) (*models.SL, error) {
	targetSL, err := s.validateSLExists(req.ID)
	if err != nil {
		return nil, err
	}
	return targetSL, nil
}

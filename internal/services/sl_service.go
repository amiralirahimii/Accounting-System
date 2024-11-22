package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/sl"
	"errors"

	"gorm.io/gorm"
)

func (s *SLService) applySLCreation(req *sl.InsertRequest) (*dtos.SLDto, error) {
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
		return err
	}
	return nil
}

func (s *SLService) applySLUpdate(req *sl.UpdateRequest, targetSL *models.SL) (*dtos.SLDto, error) {
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

func (s *SLService) validateSLExists(id int) (*models.SL, error) {
	var targetSL models.SL
	if err := s.db.Where("id = ?", id).First(&targetSL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrSLNotFound
		}
		return nil, err
	}
	return &targetSL, nil
}

func (s *SLService) validateVersion(reqVersion int, targetVersion int) error {
	if reqVersion != targetVersion {
		return constants.ErrVersionOutdated
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
		return err
	}
	return nil
}

func (s *SLService) applySLDeletion(targetSL *models.SL) error {
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

func (s *SLService) validateSLGetRequest(req *sl.GetRequest) (*models.SL, error) {
	targetSL, err := s.validateSLExists(req.ID)
	if err != nil {
		return nil, err
	}
	return targetSL, nil
}

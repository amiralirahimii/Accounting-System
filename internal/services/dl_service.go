package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"errors"

	"gorm.io/gorm"
)

type DLService struct {
	db *gorm.DB
}

func (s *DLService) InitService(db *gorm.DB) {
	s.db = db
}

func (s *DLService) applyDLCreation(req *dl.InsertRequest) (*dtos.DLDto, error) {
	dl := models.DL{
		Code:       req.Code,
		Title:      req.Title,
		RowVersion: 0,
	}

	if err := s.db.Create(&dl).Error; err != nil {
		return nil, err
	}

	return mappers.ToDLDto(&dl), nil
}

func (s *DLService) validateDLInsertRequest(req *dl.InsertRequest) error {
	if err := s.validateCodeAndTitleLength(req.Code, req.Title); err != nil {
		return err
	}
	if err := s.validateCodeAndTitleUnique(req.Code, req.Title); err != nil {
		return err
	}
	return nil
}

func (s *DLService) validateCodeAndTitleUnique(code string, title string) error {
	var existingDL models.DL
	if err := s.db.Where("code = ? OR title = ?", code, title).First(&existingDL).Error; err == nil {
		if existingDL.Code == code {
			return constants.ErrCodeAlreadyExists
		}
		if existingDL.Title == title {
			return constants.ErrTitleAlreadyExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (s *DLService) validateCodeAndTitleLength(code string, title string) error {
	if code == "" || len(code) > 64 {
		return constants.ErrCodeEmptyOrTooLong
	}
	if title == "" || len(title) > 64 {
		return constants.ErrTitleEmptyOrTooLong
	}
	return nil
}

func (s *DLService) applyDLUpdate(req *dl.UpdateRequest, targetDL *models.DL) (*dtos.DLDto, error) {
	targetDL.Code = req.Code
	targetDL.Title = req.Title
	targetDL.RowVersion++
	if err := s.db.Save(targetDL).Error; err != nil {
		return nil, err
	}

	return mappers.ToDLDto(targetDL), nil
}

func (s *DLService) validateDLUpdateRequest(req *dl.UpdateRequest) (*models.DL, error) {
	if err := s.validateCodeAndTitleLength(req.Code, req.Title); err != nil {
		return nil, err
	}
	targetDL, err := s.validateDLExists(req.ID)
	if err != nil {
		return nil, err
	}
	if err := s.validateVersion(req.Version, targetDL.RowVersion); err != nil {
		return nil, err
	}
	if err := s.validateCodeAndTitleUniqueWithDifferentId(req.Code, req.Title, req.ID); err != nil {
		return nil, err
	}
	return targetDL, nil
}

func (s *DLService) validateCodeAndTitleUniqueWithDifferentId(code string, title string, id int) error {
	var existingDL models.DL
	if err := s.db.Where("(code = ? OR title = ?) AND id != ?", code, title, id).First(&existingDL).Error; err == nil {
		if existingDL.Code == code {
			return constants.ErrCodeAlreadyExists
		}
		if existingDL.Title == title {
			return constants.ErrTitleAlreadyExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (s *DLService) validateVersion(reqVersion int, targetVersion int) error {
	if reqVersion != targetVersion {
		return constants.ErrVersionOutdated
	}
	return nil
}

func (s *DLService) validateDLExists(id int) (*models.DL, error) {
	var targetDL models.DL
	if err := s.db.Where("id = ?", id).First(&targetDL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrDLNotFound
		}
		return nil, err
	}
	return &targetDL, nil
}

func (s *DLService) applyDLDeletion(targetDL *models.DL) error {
	if err := s.db.Delete(&targetDL).Error; err != nil {
		return err
	}
	return nil
}

func (s *DLService) validateDLDeleteRequest(req *dl.DeleteRequest) (*models.DL, error) {
	targetDL, err := s.validateDLExists(req.ID)
	if err != nil {
		return nil, err
	}
	if err := s.validateVersion(req.Version, targetDL.RowVersion); err != nil {
		return nil, err
	}
	if err := s.validateDLHasNoReferences(req.ID); err != nil {
		return nil, err
	}
	return targetDL, nil
}

func (s *DLService) validateDLHasNoReferences(id int) error {
	var VoucherItemRefrencingThisDL models.VoucherItem
	if err := s.db.Where("dl_id = ?", id).First(&VoucherItemRefrencingThisDL).Error; err == nil {
		return constants.ErrThereIsRefrenceToDL
	}
	return nil
}

func (s *DLService) validateDLGetRequest(req *dl.GetRequest) (*models.DL, error) {
	targetDL, err := s.validateDLExists(req.ID)
	if err != nil {
		return nil, err
	}
	return targetDL, nil
}

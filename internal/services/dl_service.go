package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"errors"
	"log"

	"gorm.io/gorm"
)

type DLService struct {
	db *gorm.DB
}

func (s *DLService) CreateDL(req *dl.InsertRequest) (*dtos.DLDto, error) {
	if err := s.validateDLInsertRequest(req); err != nil {
		return nil, err
	}

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
		log.Printf("unexpected error while checking for duplicates: %v", err)
		return constants.ErrUnexpectedError
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

func (s *DLService) UpdateDL(req *dl.UpdateRequest) (*dtos.DLDto, error) {
	targetDL, err := s.validateDLUpdateRequest(req)
	if err != nil {
		return nil, err
	}

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
		log.Printf("unexpected error while checking for duplicates: %v", err)
		return constants.ErrUnexpectedError
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
		log.Printf("unexpected error while finding DL: %v", err)
		return nil, constants.ErrUnexpectedError
	}
	return &targetDL, nil
}

func (s *DLService) DeleteDL(req *dl.DeleteRequest) error {
	targetDL, err := s.validateDLDeleteRequest(req)
	if err != nil {
		return err
	}

	if err := s.db.Delete(&targetDL).Error; err != nil {
		log.Printf("unexpected error while deleting DL: %v", err)
		return constants.ErrUnexpectedError
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

func (s *DLService) GetDL(req *dl.GetRequest) (*dtos.DLDto, error) {
	targetDL, err := s.validateDLGetRequest(req)
	if err != nil {
		return nil, err
	}

	return mappers.ToDLDto(targetDL), nil
}

func (s *DLService) validateDLGetRequest(req *dl.GetRequest) (*models.DL, error) {
	targetDL, err := s.validateDLExists(req.ID)
	if err != nil {
		return nil, err
	}
	return targetDL, nil
}

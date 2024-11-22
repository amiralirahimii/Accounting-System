package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/requests/dl"
	"log"

	"gorm.io/gorm"
)

type DLService struct {
	db *gorm.DB
}

func (s *DLService) InitService(db *gorm.DB) {
	s.db = db
}

func (s *DLService) CreateDL(req *dl.InsertRequest) (*dtos.DLDto, error) {
	if err := s.validateDLInsertRequest(req); err != nil {
		return nil, err
	}

	dl, err := s.applyDLCreation(req)
	if err != nil {
		log.Printf("unexpected error while creating DL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return dl, nil
}

func (s *DLService) UpdateDL(req *dl.UpdateRequest) (*dtos.DLDto, error) {
	targetDL, err := s.validateDLUpdateRequest(req)
	if err != nil {
		return nil, err
	}

	dl, err := s.applyDLUpdate(req, targetDL)
	if err != nil {
		log.Printf("unexpected error while updating DL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return dl, nil
}

func (s *DLService) DeleteDL(req *dl.DeleteRequest) error {
	targetDL, err := s.validateDLDeleteRequest(req)
	if err != nil {
		return err
	}

	if err := s.applyDLDeletion(targetDL); err != nil {
		log.Printf("unexpected error while deleting DL: %v", err)
		return constants.ErrUnexpectedError
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

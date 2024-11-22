package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/requests/sl"
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

	slDto, err := s.applySLCreation(req)
	if err != nil {
		log.Printf("unexpected error while creating SL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return slDto, nil
}

func (s *SLService) UpdateSL(req *sl.UpdateRequest) (*dtos.SLDto, error) {
	targetSL, err := s.validateSLUpdateRequest(req)
	if err != nil {
		return nil, err
	}

	slDto, err := s.applySLUpdate(req, targetSL)
	if err != nil {
		log.Printf("unexpected error while updating SL: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return slDto, nil
}

func (s *SLService) DeleteSL(req *sl.DeleteRequest) error {
	targetSL, err := s.validateSLDeleteRequest(req)
	if err != nil {
		return err
	}

	if err := s.applySLDeletion(targetSL); err != nil {
		log.Printf("unexpected error while deleting SL: %v", err)
		return constants.ErrUnexpectedError
	}

	return nil
}

func (s *SLService) GetSL(req *sl.GetRequest) (*dtos.SLDto, error) {
	targetSL, err := s.validateSLGetRequest(req)
	if err != nil {
		return nil, err
	}

	return mappers.ToSlDto(targetSL), nil
}

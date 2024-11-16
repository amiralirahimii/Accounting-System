package services

import (
	"accountingsystem/db"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"errors"
)

type DLService struct{}

func (s *DLService) CreateDL(req *dl.InsertRequest) (*models.DL, error) {
	if req.Code == "" || req.Title == "" {
		return nil, errors.New("code and title cannot be empty")
	}
	if len(req.Code) > 64 || len(req.Title) > 64 {
		return nil, errors.New("code and title must be 64 characters or less")
	}

	dl := &models.DL{
		Code:       req.Code,
		Title:      req.Title,
		RowVersion: 0,
	}

	if err := db.DB.Create(dl).Error; err != nil {
		return nil, err
	}

	return dl, nil
}

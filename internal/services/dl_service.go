package services

import (
	"accountingsystem/config"
	"accountingsystem/internal/models"
	"errors"
)

// DLService handles operations related to the DL entity
type DLService struct{}

// CreateDLRequest represents the input data for creating a DL
type CreateDLRequest struct {
	Code  string
	Title string
}

// CreateDL creates a new DL entity
func (s *DLService) CreateDL(req *CreateDLRequest) (*models.DL, error) {
	// Validation
	if req.Code == "" || req.Title == "" {
		return nil, errors.New("code and title cannot be empty")
	}
	if len(req.Code) > 64 || len(req.Title) > 64 {
		return nil, errors.New("code and title must be 64 characters or less")
	}

	// Create the DL object
	dl := &models.DL{
		Code:       req.Code,
		Title:      req.Title,
		RowVersion: 0, // Initial row version
	}

	// Save to database
	if err := config.DB.Create(dl).Error; err != nil {
		return nil, err
	}

	return dl, nil
}

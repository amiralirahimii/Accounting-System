package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/sl"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomSL(s *SLService) (*models.SL, error) {
	randomCode := "SL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &sl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
		HasDL: false,
	}
	return s.CreateSL(req)
}

func Test_CreateSL_Succeeds_WithValidRequest(t *testing.T) {
	service := SLService{}
	randomCode := "SL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)

	req := &sl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
		HasDL: true,
	}

	sl, err := service.CreateSL(req)

	require.Nil(t, err)
	assert.Equal(t, sl.Code, req.Code)
	assert.Equal(t, sl.Title, req.Title)
	assert.Equal(t, sl.HasDL, req.HasDL)
}

func Test_CreateSL_ReturnsErrCodeEmptyOrTooLong_WithEmptyCode(t *testing.T) {
	service := SLService{}
	randomTitle := "Test" + generateRandomString(20)

	req := &sl.InsertRequest{
		Code:  "",
		Title: randomTitle,
		HasDL: true,
	}

	sl, err := service.CreateSL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, sl)
}

func Test_CreateSL_ReturnsErrCodeEmptyOrTooLong_WithTooLongCode(t *testing.T) {
	service := SLService{}
	randomCode := generateRandomString(65)
	randomTitle := "Test" + generateRandomString(20)

	req := &sl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
		HasDL: false,
	}

	sl, err := service.CreateSL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, sl)
}

func Test_CreateSL_ReturnsErrTitleEmptyOrTooLong_WithEmptyTitle(t *testing.T) {
	service := SLService{}
	randomCode := "SL" + generateRandomString(20)

	req := &sl.InsertRequest{
		Code:  randomCode,
		Title: "",
		HasDL: false,
	}

	sl, err := service.CreateSL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, sl)
}

func Test_CreateSL_ReturnsErrTitleEmptyOrTooLong_WithTooLongTitle(t *testing.T) {
	service := SLService{}
	randomCode := "SL" + generateRandomString(20)
	randomTitle := generateRandomString(65)

	req := &sl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
		HasDL: true,
	}

	sl, err := service.CreateSL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, sl)
}

func Test_CreateSL_ReturnsErrCodeAlreadyExists_WithExistingCode(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	duplicateReq := &sl.InsertRequest{
		Code:  createdSL.Code,
		Title: "AnotherTitle" + generateRandomString(20),
		HasDL: true,
	}

	sl, err := service.CreateSL(duplicateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeAlreadyExists)
	assert.Nil(t, sl)
}

func Test_CreateSL_ReturnsErrTitleAlreadyExists_WithExistingTitle(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	duplicateReq := &sl.InsertRequest{
		Code:  "NewCode" + generateRandomString(20),
		Title: createdSL.Title,
		HasDL: true,
	}

	sl, err := service.CreateSL(duplicateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleAlreadyExists)
	assert.Nil(t, sl)
}

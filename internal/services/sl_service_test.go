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

func Test_UpdateSL_Succeeds_WithValidRequest(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomCode := "SL" + generateRandomString(20)
	newRandomTitle := "UpdatedTitle" + generateRandomString(20)
	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		HasDL:   true,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.Nil(t, err)
	assert.Equal(t, updatedSL.Code, updateReq.Code)
	assert.Equal(t, updatedSL.Title, updateReq.Title)
	assert.Equal(t, updatedSL.HasDL, updateReq.HasDL)
}

func Test_UpdateSL_ReturnsErrCodeEmptyOrTooLong_WithEmptyCode(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomTitle := "UpdatedTitle" + generateRandomString(20)
	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    "",
		Title:   newRandomTitle,
		HasDL:   false,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedSL)
}

func Test_UpdateSL_ReturnsErrCodeEmptyOrTooLong_WithTooLongCode(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomCode := generateRandomString(65)
	newRandomTitle := "UpdatedTitle" + generateRandomString(20)
	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		HasDL:   true,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedSL)
}

func Test_UpdateSL_ReturnsErrTitleEmptyOrTooLong_WithEmptyTitle(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomCode := "SL" + generateRandomString(20)
	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode,
		Title:   "",
		HasDL:   false,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedSL)
}

func Test_UpdateSL_ReturnsErrTitleEmptyOrTooLong_WithTooLongTitle(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomCode := "SL" + generateRandomString(20)
	newRandomTitle := generateRandomString(65)
	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		HasDL:   true,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedSL)
}

func Test_UpdateSL_ReturnsErrSLNotFound_WithNonExistingID(t *testing.T) {
	service := SLService{}
	newRandomCode := "SL" + generateRandomString(20)
	newRandomTitle := "NewTitle" + generateRandomString(20)

	updateReq := &sl.UpdateRequest{
		ID:    generateRandomInt64(),
		Code:  newRandomCode,
		Title: newRandomTitle,
		HasDL: true,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrSLNotFound)
	assert.Nil(t, updatedSL)
}

func Test_UpdateSL_ReturnsErrVersionOutdated_WithOutdatedVersion(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomCode := "SL" + generateRandomString(20)
	newRandomTitle := "NewTitle" + generateRandomString(20)

	updateReq1 := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		HasDL:   true,
		Version: createdSL.RowVersion,
	}

	_, err = service.UpdateSL(updateReq1)
	require.Nil(t, err)

	newRandomCode2 := "SL" + generateRandomString(20)
	newRandomTitle2 := "NewTitle" + generateRandomString(20)

	updateReq2 := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode2,
		Title:   newRandomTitle2,
		HasDL:   false,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq2)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
	assert.Nil(t, updatedSL)
}

func Test_UpdateSL_ReturnsErrCodeAlreadyExists_WithExistingCode(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	duplicateSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomTitle := "NewTitle" + generateRandomString(20)

	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    duplicateSL.Code,
		Title:   newRandomTitle,
		HasDL:   true,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeAlreadyExists)
	assert.Nil(t, updatedSL)
}

func Test_UpdateSL_ReturnsErrTitleAlreadyExists_WithExistingTitle(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	duplicateSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomCode := "SL" + generateRandomString(20)

	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode,
		Title:   duplicateSL.Title,
		HasDL:   true,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := service.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleAlreadyExists)
	assert.Nil(t, updatedSL)
}

func Test_DeleteSL_Succeeds_WithValidRequest(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	deleteReq := &sl.DeleteRequest{
		ID:      createdSL.ID,
		Version: createdSL.RowVersion,
	}

	err = service.DeleteSL(deleteReq)

	require.Nil(t, err)
}

func Test_DeleteSL_ReturnsErrSLNotFound_WithNonExistingID(t *testing.T) {
	service := SLService{}
	deleteReq := &sl.DeleteRequest{
		ID: generateRandomInt64(),
	}

	err := service.DeleteSL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrSLNotFound)
}

func Test_DeleteSL_ReturnsErrSLNotFound_WithDeletedID(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	deleteReq := &sl.DeleteRequest{
		ID:      createdSL.ID,
		Version: createdSL.RowVersion,
	}

	err = service.DeleteSL(deleteReq)
	require.Nil(t, err)

	err = service.DeleteSL(deleteReq)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrSLNotFound)
}

func Test_DeleteSL_ReturnsErrVersionOutdated_WithOutdatedVersion(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service)
	require.Nil(t, err)

	newRandomCode := "SL" + generateRandomString(20)
	newRandomTitle := "NewTitle" + generateRandomString(20)
	newHasDL := !createdSL.HasDL
	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		HasDL:   newHasDL,
		Version: createdSL.RowVersion,
	}

	_, err = service.UpdateSL(updateReq)
	require.Nil(t, err)

	deleteReq := &sl.DeleteRequest{
		ID:      createdSL.ID,
		Version: createdSL.RowVersion,
	}

	err = service.DeleteSL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
}

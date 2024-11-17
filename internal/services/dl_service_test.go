package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateDL_Succeeds_WithValidRequest(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.Nil(t, err)
	assert.Equal(t, dl.Code, req.Code)
	assert.Equal(t, dl.Title, dl.Title)
}

func Test_CreateDL_ReturnsErrCodeEmptyOrTooLong_WithEmptyCode(t *testing.T) {
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  "",
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrCodeEmptyOrTooLong_WithTooLongCode(t *testing.T) {
	randomCode := generateRandomString(65)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrTitleEmptyOrTooLong_WithEmptyTitle(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: "",
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrTitleEmptyOrTooLong_WithTooLongTitle(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := generateRandomString(65)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrCodeAlreadyExists_WithExistingCode(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	validReq := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	randomTitleNotExisting := "Test" + generateRandomString(20)
	reqWithExistingCode := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitleNotExisting,
	}

	service := DLService{}
	_, err := service.CreateDL(validReq)
	require.Nil(t, err)

	dl, err := service.CreateDL(reqWithExistingCode)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeAlreadyExists)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrTitleAlreadyExists_WithExistingTitle(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	validReq := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	randomCodeNotExisting := "DL" + generateRandomString(20)
	reqWithExistingTitle := &dl.InsertRequest{
		Code:  randomCodeNotExisting,
		Title: randomTitle,
	}

	service := DLService{}
	_, err := service.CreateDL(validReq)
	require.Nil(t, err)

	dl, err := service.CreateDL(reqWithExistingTitle)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleAlreadyExists)
	assert.Nil(t, dl)
}

func createRandomDL(s *DLService) (*models.DL, error) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	createDL, err := s.CreateDL(req)
	return createDL, err
}

func Test_UpdateDL_Succeeds_WithValidRequest(t *testing.T) {
	service := DLService{}
	createdDL, err := createRandomDL(&service)
	require.Nil(t, err)

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.Nil(t, err)
	assert.Equal(t, updatedDL.Code, updateReq.Code)
	assert.Equal(t, updatedDL.Title, updateReq.Title)
}

func Test_UpdateDL_ReturnsErrCodeEmptyOrTooLong_WithEmptyCode(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    "",
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrCodeEmptyOrTooLong_WithTooLongCode(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := generateRandomString(65)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleEmptyOrTooLong_WithEmptyTitle(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := "DL" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   "",
		Version: createdDL.RowVersion,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleEmptyOrTooLong_WithTooLongTitle(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := generateRandomString(65)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrDLNotFound_WithNonExistingID(t *testing.T) {
	service := DLService{}
	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	newID := generateRandomInt64()

	updateReq := &dl.UpdateRequest{
		ID:    newID,
		Code:  newRandomCode,
		Title: newRandomTitle,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrVersionOutdated_WithOutdatedVersion(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	newRandomCode2 := "DL" + generateRandomString(20)
	newRandomTitle2 := "Test" + generateRandomString(20)
	updateReqWithOutdatedVersion := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode2,
		Title:   newRandomTitle2,
		Version: createdDL.RowVersion,
	}

	_, err := service.UpdateDL(updateReq)
	require.Nil(t, err)

	updatedDL, err := service.UpdateDL(updateReqWithOutdatedVersion)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrCodeAlreadyExists_WithExistingCode(t *testing.T) {
	service := DLService{}
	createdDL1, _ := createRandomDL(&service)
	createdDL2, _ := createRandomDL(&service)

	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL1.ID,
		Code:    createdDL2.Code,
		Title:   newRandomTitle,
		Version: createdDL1.RowVersion,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeAlreadyExists)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleAlreadyExists_WithExistingTitle(t *testing.T) {
	service := DLService{}
	createdDL1, _ := createRandomDL(&service)
	createdDL2, _ := createRandomDL(&service)

	newRandomCode := "DL" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL1.ID,
		Code:    newRandomCode,
		Title:   createdDL2.Title,
		Version: createdDL1.RowVersion,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleAlreadyExists)
	assert.Nil(t, updatedDL)
}

func Test_DeleteDL_Succeeds_WithValidRequest(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	deleteReq := &dl.DeleteRequest{
		ID:      createdDL.ID,
		Version: createdDL.RowVersion,
	}

	err := service.DeleteDL(deleteReq)

	require.Nil(t, err)
}

func Test_DeleteDL_ReturnsErrDLNotFound_WithNonExistingID(t *testing.T) {
	service := DLService{}
	newID := generateRandomInt64()

	deleteReq := &dl.DeleteRequest{
		ID: newID,
	}

	err := service.DeleteDL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
}

func Test_DeleteDL_ReturnsErrDLNotFound_WithDeletedID(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	deleteReq := &dl.DeleteRequest{
		ID:      createdDL.ID,
		Version: createdDL.RowVersion,
	}

	err := service.DeleteDL(deleteReq)

	require.Nil(t, err)

	err = service.DeleteDL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
}

func Test_DeleteDL_ReturnsErrVersionOutdated_WithOutdatedVersion(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	_, err := service.UpdateDL(updateReq)
	require.Nil(t, err)

	deleteReq := &dl.DeleteRequest{
		ID:      createdDL.ID,
		Version: createdDL.RowVersion,
	}

	err = service.DeleteDL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
}

func Test_GetDLByID_Succeeds_WithValidID(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	getReq := &dl.GetRequest{
		ID: createdDL.ID,
	}

	foundDL, err := service.GetDL(getReq)

	require.Nil(t, err)
	assert.Equal(t, createdDL.ID, foundDL.ID)
	assert.Equal(t, createdDL.Code, foundDL.Code)
	assert.Equal(t, createdDL.Title, foundDL.Title)
	assert.Equal(t, createdDL.RowVersion, foundDL.RowVersion)
}

func Test_GetDLByID_ReturnsErrDLNotFound_WithNonExistingID(t *testing.T) {
	service := DLService{}
	newID := generateRandomInt64()

	getReq := &dl.GetRequest{
		ID: newID,
	}

	foundDL, err := service.GetDL(getReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
	assert.Nil(t, foundDL)
}

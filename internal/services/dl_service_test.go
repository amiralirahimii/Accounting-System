package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/requests/dl"
	"accountingsystem/internal/requests/sl"
	"accountingsystem/internal/requests/voucher"
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

	dl, err := dlService.CreateDL(req)

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

	dl, err := dlService.CreateDL(req)

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

	dl, err := dlService.CreateDL(req)

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

	dl, err := dlService.CreateDL(req)

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

	dl, err := dlService.CreateDL(req)

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

	_, err := dlService.CreateDL(validReq)
	require.Nil(t, err)

	randomTitleNotExisting := "Test" + generateRandomString(20)
	reqWithExistingCode := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitleNotExisting,
	}

	dl, err := dlService.CreateDL(reqWithExistingCode)

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

	_, err := dlService.CreateDL(validReq)
	require.Nil(t, err)

	randomCodeNotExisting := "DL" + generateRandomString(20)
	reqWithExistingTitle := &dl.InsertRequest{
		Code:  randomCodeNotExisting,
		Title: randomTitle,
	}

	dl, err := dlService.CreateDL(reqWithExistingTitle)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleAlreadyExists)
	assert.Nil(t, dl)
}

func Test_UpdateDL_Succeeds_WithValidRequest(t *testing.T) {
	createdDL, err := createRandomDL()
	require.Nil(t, err)

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.Nil(t, err)
	assert.Equal(t, updatedDL.Code, updateReq.Code)
	assert.Equal(t, updatedDL.Title, updateReq.Title)
}

func Test_UpdateDL_ReturnsErrCodeEmptyOrTooLong_WithEmptyCode(t *testing.T) {
	createdDL, _ := createRandomDL()

	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    "",
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrCodeEmptyOrTooLong_WithTooLongCode(t *testing.T) {
	createdDL, _ := createRandomDL()

	newRandomCode := generateRandomString(65)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleEmptyOrTooLong_WithEmptyTitle(t *testing.T) {
	createdDL, _ := createRandomDL()

	newRandomCode := "DL" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   "",
		Version: createdDL.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleEmptyOrTooLong_WithTooLongTitle(t *testing.T) {
	createdDL, _ := createRandomDL()

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := generateRandomString(65)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrDLNotFound_WithNonExistingID(t *testing.T) {
	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	newID := generateRandomInt64()

	updateReq := &dl.UpdateRequest{
		ID:    newID,
		Code:  newRandomCode,
		Title: newRandomTitle,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrVersionOutdated_WithOutdatedVersion(t *testing.T) {
	createdDL, _ := createRandomDL()

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	_, err := dlService.UpdateDL(updateReq)
	require.Nil(t, err)

	newRandomCode2 := "DL" + generateRandomString(20)
	newRandomTitle2 := "Test" + generateRandomString(20)
	updateReqWithOutdatedVersion := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode2,
		Title:   newRandomTitle2,
		Version: createdDL.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReqWithOutdatedVersion)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrCodeAlreadyExists_WithExistingCode(t *testing.T) {
	createdDL1, _ := createRandomDL()
	createdDL2, _ := createRandomDL()

	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL1.ID,
		Code:    createdDL2.Code,
		Title:   newRandomTitle,
		Version: createdDL1.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeAlreadyExists)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleAlreadyExists_WithExistingTitle(t *testing.T) {
	createdDL1, _ := createRandomDL()
	createdDL2, _ := createRandomDL()

	newRandomCode := "DL" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL1.ID,
		Code:    newRandomCode,
		Title:   createdDL2.Title,
		Version: createdDL1.RowVersion,
	}

	updatedDL, err := dlService.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleAlreadyExists)
	assert.Nil(t, updatedDL)
}

func Test_DeleteDL_Succeeds_WithValidRequest(t *testing.T) {
	createdDL, _ := createRandomDL()

	deleteReq := &dl.DeleteRequest{
		ID:      createdDL.ID,
		Version: createdDL.RowVersion,
	}

	err := dlService.DeleteDL(deleteReq)

	require.Nil(t, err)
}

func Test_DeleteDL_ReturnsErrDLNotFound_WithNonExistingID(t *testing.T) {
	newID := generateRandomInt64()

	deleteReq := &dl.DeleteRequest{
		ID: newID,
	}

	err := dlService.DeleteDL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
}

func Test_DeleteDL_ReturnsErrDLNotFound_WithDeletedID(t *testing.T) {
	createdDL, _ := createRandomDL()

	deleteReq := &dl.DeleteRequest{
		ID:      createdDL.ID,
		Version: createdDL.RowVersion,
	}

	err := dlService.DeleteDL(deleteReq)

	require.Nil(t, err)

	err = dlService.DeleteDL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
}

func Test_DeleteDL_ReturnsErrVersionOutdated_WithOutdatedVersion(t *testing.T) {
	createdDL, _ := createRandomDL()

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:      createdDL.ID,
		Code:    newRandomCode,
		Title:   newRandomTitle,
		Version: createdDL.RowVersion,
	}

	_, err := dlService.UpdateDL(updateReq)
	require.Nil(t, err)

	deleteReq := &dl.DeleteRequest{
		ID:      createdDL.ID,
		Version: createdDL.RowVersion,
	}

	err = dlService.DeleteDL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
}

func Test_DeleteDL_ReturnsErrThereIsRefrenceToDL_WithExistingReference(t *testing.T) {
	slService := SLService{}
	slWithDL, err := slService.CreateSL(&sl.InsertRequest{
		Code:  "SL" + generateRandomString(20),
		Title: "SLWithDL" + generateRandomString(20),
		HasDL: true,
	})
	require.Nil(t, err)

	createdDL, err := createRandomDL()
	require.Nil(t, err)

	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   slWithDL.ID,
			DLID:   &createdDL.ID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   slWithDL.ID,
			DLID:   &createdDL.ID,
			Debit:  0,
			Credit: 100,
		},
	}
	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}
	_, err = voucherService.CreateVoucher(req)
	require.Nil(t, err)

	deleteReq := &dl.DeleteRequest{
		ID:      createdDL.ID,
		Version: createdDL.RowVersion,
	}
	err = dlService.DeleteDL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrThereIsRefrenceToDL)
}

func Test_GetDLByID_Succeeds_WithValidID(t *testing.T) {
	createdDL, _ := createRandomDL()

	getReq := &dl.GetRequest{
		ID: createdDL.ID,
	}

	foundDL, err := dlService.GetDL(getReq)

	require.Nil(t, err)
	assert.Equal(t, createdDL.ID, foundDL.ID)
	assert.Equal(t, createdDL.Code, foundDL.Code)
	assert.Equal(t, createdDL.Title, foundDL.Title)
	assert.Equal(t, createdDL.RowVersion, foundDL.RowVersion)
}

func Test_GetDLByID_ReturnsErrDLNotFound_WithNonExistingID(t *testing.T) {
	newID := generateRandomInt64()

	getReq := &dl.GetRequest{
		ID: newID,
	}

	foundDL, err := dlService.GetDL(getReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
	assert.Nil(t, foundDL)
}

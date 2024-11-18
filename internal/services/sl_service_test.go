package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/requests/sl"
	"accountingsystem/internal/requests/voucher"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
	require.Nil(t, err)

	duplicateSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
	require.Nil(t, err)

	duplicateSL, err := createRandomSL(&service, true)
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

func Test_UpdateSL_ReturnsErrThereIsReferenceToSL_WithReferencedSL(t *testing.T) {
	slService := SLService{}
	createdSL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	dlService := DLService{}
	createdDL, err := createRandomDL(&dlService)
	require.Nil(t, err)

	voucherService := VoucherService{}
	voucherItems := []voucher.VoucherItemInsertRequest{
		{
			SLID:   createdSL.ID,
			DLID:   &createdDL.ID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   createdSL.ID,
			DLID:   &createdDL.ID,
			Debit:  0,
			Credit: 100,
		},
	}
	voucherReq := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: voucherItems,
	}
	_, err = voucherService.CreateVoucher(voucherReq)
	require.Nil(t, err)

	updateReq := &sl.UpdateRequest{
		ID:      createdSL.ID,
		Code:    "UpdatedCode" + generateRandomString(20),
		Title:   "UpdatedTitle" + generateRandomString(20),
		HasDL:   createdSL.HasDL,
		Version: createdSL.RowVersion,
	}

	updatedSL, err := slService.UpdateSL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrThereIsRefrenceToSL)
	assert.Nil(t, updatedSL)
}

func Test_DeleteSL_Succeeds_WithValidRequest(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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
	createdSL, err := createRandomSL(&service, true)
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

func Test_GetSL_Succeeds_WithValidID(t *testing.T) {
	service := SLService{}
	createdSL, err := createRandomSL(&service, true)
	require.Nil(t, err)

	getReq := &sl.GetRequest{
		ID: createdSL.ID,
	}

	sl, err := service.GetSL(getReq)

	require.Nil(t, err)
	assert.Equal(t, sl.ID, createdSL.ID)
	assert.Equal(t, sl.Code, createdSL.Code)
	assert.Equal(t, sl.Title, createdSL.Title)
	assert.Equal(t, sl.HasDL, createdSL.HasDL)
}

func Test_GetSL_ReturnsErrSLNotFound_WithNonExistingID(t *testing.T) {
	service := SLService{}
	getReq := &sl.GetRequest{
		ID: generateRandomInt64(),
	}

	sl, err := service.GetSL(getReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrSLNotFound)
	assert.Nil(t, sl)
}

func Test_DeleteSL_ReturnsErrThereIsReferenceToSL_WithReferencedSL(t *testing.T) {
	slService := SLService{}
	createdSL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	dlService := DLService{}
	createdDL, err := createRandomDL(&dlService)
	require.Nil(t, err)

	voucherService := VoucherService{}
	voucherItems := []voucher.VoucherItemInsertRequest{
		{
			SLID:   createdSL.ID,
			DLID:   &createdDL.ID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   createdSL.ID,
			DLID:   &createdDL.ID,
			Debit:  0,
			Credit: 100,
		},
	}
	voucherReq := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: voucherItems,
	}
	_, err = voucherService.CreateVoucher(voucherReq)
	require.Nil(t, err)

	deleteReq := &sl.DeleteRequest{
		ID:      createdSL.ID,
		Version: createdSL.RowVersion,
	}

	err = slService.DeleteSL(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrThereIsRefrenceToSL)
}

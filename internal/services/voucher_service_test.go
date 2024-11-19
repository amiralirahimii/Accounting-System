package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/requests/voucher"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var voucherService = VoucherService{}
var slService = SLService{}
var dlService = DLService{}

func Test_CreateVoucher_Succeeds_ReferencingDLAndNonReferencingDLVoucherItems(t *testing.T) {
	slWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	slWithoutDL, err := createRandomSL(&slService, false)
	require.Nil(t, err)

	dl, err := createRandomDL(&dlService)
	require.Nil(t, err)

	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   slWithDL.ID,
			DLID:   &dl.ID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   slWithoutDL.ID,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)

	require.Nil(t, err)
	assert.Equal(t, req.Number, voucher.Number)
}

func Test_CreateVoucher_ReturnsErrNumberEmptyOrTooLong_WithTooShortNumber(t *testing.T) {
	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   1,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   2,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       "",
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrNumberEmptyOrTooLong)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrNumberEmptyOrTooLong_WithTooLongNumber(t *testing.T) {
	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   1,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   2,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(65),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrNumberEmptyOrTooLong)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrVoucherItemsCountOutOfRange_WithLessThanTwoVoucherItems(t *testing.T) {
	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   1,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherItemsCountOutOfRange)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrVoucherItemsCountOutOfRange_WithMoreThan500VoucherItems(t *testing.T) {
	var items []voucher.VoucherItemInsertDetail
	for i := 0; i < 501; i++ {
		items = append(items, voucher.VoucherItemInsertDetail{
			SLID:   1,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		})
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherItemsCountOutOfRange)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrDebitOrCreditInvalid_BothDebitAndCreditNonZero(t *testing.T) {
	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   1,
			DLID:   nil,
			Debit:  100,
			Credit: 50,
		},
		{
			SLID:   2,
			DLID:   nil,
			Debit:  0,
			Credit: 50,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDebitOrCreditInvalid)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrDebitOrCreditInvalid_BothDebitAndCreditZero(t *testing.T) {
	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   1,
			DLID:   nil,
			Debit:  0,
			Credit: 0,
		},
		{
			SLID:   2,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
		{
			SLID:   3,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDebitOrCreditInvalid)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrDebitCreditMismatch_WithMismatchedDebitCreditSum(t *testing.T) {
	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   1,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   2,
			DLID:   nil,
			Debit:  0,
			Credit: 50,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDebitCreditMismatch)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrVoucherNumberExists_WithExistingVoucherNumber(t *testing.T) {
	slWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	slWithoutDL, err := createRandomSL(&slService, false)
	require.Nil(t, err)

	dl, err := createRandomDL(&dlService)
	require.Nil(t, err)

	initialItems := []voucher.VoucherItemInsertDetail{
		{
			SLID:   slWithDL.ID,
			DLID:   &dl.ID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   slWithoutDL.ID,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	initialReq := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: initialItems,
	}

	_, err = voucherService.CreateVoucher(initialReq)
	require.Nil(t, err)

	newItems := []voucher.VoucherItemInsertDetail{
		{
			SLID:   1,
			DLID:   nil,
			Debit:  200,
			Credit: 0,
		},
		{
			SLID:   2,
			DLID:   nil,
			Debit:  0,
			Credit: 200,
		},
	}

	newReq := &voucher.InsertRequest{
		Number:       initialReq.Number,
		VoucherItems: newItems,
	}

	voucher, err := voucherService.CreateVoucher(newReq)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherNumberExists)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrSLNotFound_WithInvalidSLID(t *testing.T) {
	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   generateRandomInt64(),
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   generateRandomInt64(),
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrSLNotFound)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrDLIDRequired_WhenSLRequiresDLButNoDLProvided(t *testing.T) {
	slWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   slWithDL.ID,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   slWithDL.ID,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLIDRequired)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrDLNotAllowed_WhenSLDoesNotRequireDLButDLProvided(t *testing.T) {
	slWithoutDL, err := createRandomSL(&slService, false)
	require.Nil(t, err)

	dl, err := createRandomDL(&dlService)
	require.Nil(t, err)

	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   slWithoutDL.ID,
			DLID:   &dl.ID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   slWithoutDL.ID,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotAllowed)
	assert.Nil(t, voucher)
}

func Test_CreateVoucher_ReturnsErrDLNotFound_WithInvalidDLID(t *testing.T) {
	slWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	invalidDLID := generateRandomInt64()

	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   slWithDL.ID,
			DLID:   &invalidDLID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   slWithDL.ID,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucher, err := voucherService.CreateVoucher(req)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
	assert.Nil(t, voucher)
}

func createRandomVoucher() (*dtos.VoucherDto, error) {
	slWithDL, err := createRandomSL(&slService, true)
	if err != nil {
		return nil, err
	}

	slWithoutDL, err := createRandomSL(&slService, false)
	if err != nil {
		return nil, err
	}

	dl, err := createRandomDL(&dlService)
	if err != nil {
		return nil, err
	}

	items := []voucher.VoucherItemInsertDetail{
		{
			SLID:   slWithDL.ID,
			DLID:   &dl.ID,
			Debit:  100,
			Credit: 0,
		},
		{
			SLID:   slWithoutDL.ID,
			DLID:   nil,
			Debit:  0,
			Credit: 100,
		},
	}

	req := &voucher.InsertRequest{
		Number:       generateRandomString(20),
		VoucherItems: items,
	}

	voucherDto, err := voucherService.CreateVoucher(req)
	if err != nil {
		return nil, err
	}

	return voucherDto, nil
}

func Test_UpdateVoucher_Succeeds_WithValidRequest(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	SLWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	SLWithoutDL, err := createRandomSL(&slService, false)
	require.Nil(t, err)

	DL, err := createRandomDL(&dlService)
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   SLWithDL.ID,
				DLID:   &DL.ID,
				Debit:  1000,
				Credit: 0,
			},
			{
				SLID:   SLWithoutDL.ID,
				DLID:   nil,
				Debit:  0,
				Credit: 500,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     voucherDto.VoucherItems[0].ID,
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   &voucherDto.VoucherItems[0].DLID,
				Debit:  0,
				Credit: 500,
			},
		},
		Deleted: []int{
			voucherDto.VoucherItems[1].ID,
		},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.Nil(t, err)
	assert.Equal(t, req.Number, voucher.Number)
}

func Test_UpdateVoucher_ReturnsErrVoucherNotFound_WithNonExistentVoucherID(t *testing.T) {
	insertReq, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   insertReq.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   insertReq.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  0,
				Credit: 100,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      generateRandomInt64(),
		Version: 0,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherNotFound)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrNumberEmptyOrTooLong_WithEmptyNumber(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  0,
				Credit: 100,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  "",
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrNumberEmptyOrTooLong)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrNumberEmptyOrTooLong_WithTooLongNumber(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  0,
				Credit: 100,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(65),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrNumberEmptyOrTooLong)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrVersionOutdated_WithOutdatedRequest(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	SLWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	DL, err := createRandomDL(&dlService)
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   SLWithDL.ID,
				DLID:   &DL.ID,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   SLWithDL.ID,
				DLID:   &DL.ID,
				Debit:  0,
				Credit: 100,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	_, err = voucherService.UpdateVoucher(req)

	require.Nil(t, err)

	updatedVoucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
	assert.Nil(t, updatedVoucher)
}

func Test_UpdateVoucher_ReturnsErrSLNotFound_WithInvalidSLInInserts(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   generateRandomInt64(),
				DLID:   nil,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   generateRandomInt64(),
				DLID:   nil,
				Debit:  0,
				Credit: 100,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrSLNotFound)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrDLIDRequired_WhenSLRequiresDLButNoDLProvidedInInserts(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	slWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   slWithDL.ID,
				DLID:   nil,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   slWithDL.ID,
				DLID:   nil,
				Debit:  0,
				Credit: 100,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLIDRequired)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrDLNotAllowed_WhenSLDoesNotRequireDLButDLProvidedInInserts(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	slWithoutDL, err := createRandomSL(&slService, false)
	require.Nil(t, err)

	dl, err := createRandomDL(&dlService)
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   slWithoutDL.ID,
				DLID:   &dl.ID,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   slWithoutDL.ID,
				DLID:   &dl.ID,
				Debit:  0,
				Credit: 100,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotAllowed)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrDebitOrCreditInvalid_InvalidValuesInInserts(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  100,
				Credit: 50,
			},
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  0,
				Credit: 50,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDebitOrCreditInvalid)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrSLNotFound_WithInvalidSLInUpdates(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     voucherDto.VoucherItems[0].ID,
				SLID:   generateRandomInt64(),
				DLID:   nil,
				Debit:  voucherDto.VoucherItems[0].Debit,
				Credit: voucherDto.VoucherItems[0].Credit,
			},
		},
		Inserted: []voucher.VoucherItemInsertDetail{},
		Deleted:  []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrSLNotFound)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrDLIDRequired_WhenSLRequiresDLButNoDLProvidedInUpdates(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	slWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     voucherDto.VoucherItems[0].ID,
				SLID:   slWithDL.ID,
				DLID:   nil,
				Debit:  voucherDto.VoucherItems[0].Debit,
				Credit: voucherDto.VoucherItems[0].Credit,
			},
		},
		Inserted: []voucher.VoucherItemInsertDetail{},
		Deleted:  []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLIDRequired)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrDLNotAllowed_WhenSLDoesNotRequireDLButDLProvidedInUpdates(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	slWithoutDL, err := createRandomSL(&slService, false)
	require.Nil(t, err)

	dl, err := createRandomDL(&dlService)
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     voucherDto.VoucherItems[0].ID,
				SLID:   slWithoutDL.ID,
				DLID:   &dl.ID,
				Debit:  voucherDto.VoucherItems[0].Debit,
				Credit: voucherDto.VoucherItems[0].Credit,
			},
		},
		Inserted: []voucher.VoucherItemInsertDetail{},
		Deleted:  []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotAllowed)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrDebitOrCreditInvalid_InvalidValuesInUpdates(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     voucherDto.VoucherItems[0].ID,
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   nil,
				Debit:  voucherDto.VoucherItems[0].Debit + 50,
				Credit: voucherDto.VoucherItems[0].Credit + 50,
			},
		},
		Inserted: []voucher.VoucherItemInsertDetail{},
		Deleted:  []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDebitOrCreditInvalid)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrVoucherItemNotFound_WithInvalidVoucherItemIDInUpdates(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	invalidVoucherItemID := generateRandomInt64()

	items := voucher.VoucherItemsUpdate{
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     invalidVoucherItemID,
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   &voucherDto.VoucherItems[0].DLID,
				Debit:  voucherDto.VoucherItems[0].Debit,
				Credit: voucherDto.VoucherItems[0].Credit,
			},
		},
		Inserted: []voucher.VoucherItemInsertDetail{},
		Deleted:  []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherItemNotFound)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrDebitCreditMismatch_WithUnbalancedRequest(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   &voucherDto.VoucherItems[0].DLID,
				Debit:  100,
				Credit: 0,
			},
			{
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   &voucherDto.VoucherItems[0].DLID,
				Debit:  100,
				Credit: 0,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     voucherDto.VoucherItems[0].ID,
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   &voucherDto.VoucherItems[0].DLID,
				Debit:  0,
				Credit: 500,
			},
		},
		Deleted: []int{voucherDto.VoucherItems[1].ID},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDebitCreditMismatch)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrVoucherItemsCountOutOfRange_WhenCountGoesBelowTwo(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{},
		Updated:  []voucher.VoucherItemUpdateDetail{},
		Deleted: []int{
			voucherDto.VoucherItems[0].ID,
		},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherItemsCountOutOfRange)
	assert.Nil(t, voucher)
}

func Test_UpdateVoucher_ReturnsErrVoucherItemsCountOutOfRange_WhenCountExceedsFiveHundred(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	var insertedItems []voucher.VoucherItemInsertDetail
	for i := 0; i < 498; i++ {
		insertedItems = append(insertedItems, voucher.VoucherItemInsertDetail{
			SLID:   voucherDto.VoucherItems[0].SLID,
			DLID:   nil,
			Debit:  100,
			Credit: 0,
		})
	}
	insertedItems = append(insertedItems, voucher.VoucherItemInsertDetail{
		SLID:   voucherDto.VoucherItems[0].SLID,
		DLID:   nil,
		Debit:  0,
		Credit: 50000,
	})

	items := voucher.VoucherItemsUpdate{
		Inserted: insertedItems,
		Updated:  []voucher.VoucherItemUpdateDetail{},
		Deleted:  []int{},
	}

	req := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	voucher, err := voucherService.UpdateVoucher(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherItemsCountOutOfRange)
	assert.Nil(t, voucher)
}

func Test_DeleteVoucher_Succeeds_WithValidRequest(t *testing.T) {
	createdVoucher, err := createRandomVoucher()
	require.Nil(t, err)

	deleteReq := &voucher.DeleteRequest{
		ID:      createdVoucher.ID,
		Version: createdVoucher.RowVersion,
	}

	err = voucherService.DeleteVoucher(deleteReq)

	require.Nil(t, err)

	_, err = voucherService.GetVoucher(&voucher.GetRequest{ID: createdVoucher.ID})
	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherNotFound)
}

func Test_DeleteVoucher_ReturnsErrVoucherNotFound_WithNonExistingVoucherID(t *testing.T) {
	deleteReq := &voucher.DeleteRequest{
		ID: generateRandomInt64(),
	}

	err := voucherService.DeleteVoucher(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherNotFound)
}

func Test_DeleteVoucher_ReturnsErrVersionOutdated_WithOutdatedVersion(t *testing.T) {
	voucherDto, err := createRandomVoucher()
	require.Nil(t, err)

	SLWithDL, err := createRandomSL(&slService, true)
	require.Nil(t, err)

	SLWithoutDL, err := createRandomSL(&slService, false)
	require.Nil(t, err)

	DL, err := createRandomDL(&dlService)
	require.Nil(t, err)

	items := voucher.VoucherItemsUpdate{
		Inserted: []voucher.VoucherItemInsertDetail{
			{
				SLID:   SLWithDL.ID,
				DLID:   &DL.ID,
				Debit:  1000,
				Credit: 0,
			},
			{
				SLID:   SLWithoutDL.ID,
				DLID:   nil,
				Debit:  0,
				Credit: 500,
			},
		},
		Updated: []voucher.VoucherItemUpdateDetail{
			{
				ID:     voucherDto.VoucherItems[0].ID,
				SLID:   voucherDto.VoucherItems[0].SLID,
				DLID:   &voucherDto.VoucherItems[0].DLID,
				Debit:  0,
				Credit: 500,
			},
		},
		Deleted: []int{
			voucherDto.VoucherItems[1].ID,
		},
	}

	validUpdateReq := &voucher.UpdateRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
		Number:  generateRandomString(20),
		Items:   items,
	}

	_, err = voucherService.UpdateVoucher(validUpdateReq)
	require.Nil(t, err)

	deleteReq := &voucher.DeleteRequest{
		ID:      voucherDto.ID,
		Version: voucherDto.RowVersion,
	}

	err = voucherService.DeleteVoucher(deleteReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVersionOutdated)
}

func Test_GetVoucher_Succeeds_WithValidRequest(t *testing.T) {
	createdVoucher, err := createRandomVoucher()
	require.Nil(t, err)

	getReq := &voucher.GetRequest{
		ID: createdVoucher.ID,
	}

	voucherDto, err := voucherService.GetVoucher(getReq)

	require.Nil(t, err)
	assert.Equal(t, createdVoucher.ID, voucherDto.ID)
	assert.Equal(t, createdVoucher.Number, voucherDto.Number)
}

func Test_GetVoucher_ReturnsErrVoucherNotFound_WithNonExistingVoucherID(t *testing.T) {
	newID := generateRandomInt64()

	getReq := &voucher.GetRequest{
		ID: newID,
	}

	voucherDto, err := voucherService.GetVoucher(getReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrVoucherNotFound)
	assert.Nil(t, voucherDto)
}

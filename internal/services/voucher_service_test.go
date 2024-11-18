package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"accountingsystem/internal/requests/sl"
	"accountingsystem/internal/requests/voucher"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var voucherService = VoucherService{}

func createRandomSLInDB(hasDL bool) (*models.SL, error) {
	slService := SLService{}
	return slService.CreateSL(&sl.InsertRequest{
		Code:  "SL" + generateRandomString(20),
		Title: "Title" + generateRandomString(20),
		HasDL: hasDL,
	})
}

func createRandomDLInDB() (*models.DL, error) {
	dlService := DLService{}
	return dlService.CreateDL(&dl.InsertRequest{
		Code:  "DL" + generateRandomString(20),
		Title: "Title" + generateRandomString(20),
	})
}

func Test_CreateVoucher_Succeeds_ReferencingDLAndNonReferencingDLVoucherItems(t *testing.T) {
	slWithDL, err := createRandomSLInDB(true)
	require.Nil(t, err)

	slWithoutDL, err := createRandomSLInDB(false)
	require.Nil(t, err)

	dl, err := createRandomDLInDB()
	require.Nil(t, err)

	items := []voucher.VoucherItemInsertRequest{
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
	items := []voucher.VoucherItemInsertRequest{
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
	items := []voucher.VoucherItemInsertRequest{
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

func Test_CreateVoucher_Fails_WithLessThanTwoVoucherItems(t *testing.T) {
	items := []voucher.VoucherItemInsertRequest{
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

func Test_CreateVoucher_Fails_WithMoreThan500VoucherItems(t *testing.T) {
	var items []voucher.VoucherItemInsertRequest
	for i := 0; i < 501; i++ {
		items = append(items, voucher.VoucherItemInsertRequest{
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


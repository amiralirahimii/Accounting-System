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

func Test_CreateVoucher_ReturnsErrVoucherItemsCountOutOfRange_WithLessThanTwoVoucherItems(t *testing.T) {
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

func Test_CreateVoucher_ReturnsErrVoucherItemsCountOutOfRange_WithMoreThan500VoucherItems(t *testing.T) {
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

func Test_CreateVoucher_ReturnsErrDebitOrCreditInvalid_BothDebitAndCreditNonZero(t *testing.T) {
	items := []voucher.VoucherItemInsertRequest{
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
			Credit: 100,
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
	items := []voucher.VoucherItemInsertRequest{
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
	slWithDL, err := createRandomSLInDB(true)
	require.Nil(t, err)

	slWithoutDL, err := createRandomSLInDB(false)
	require.Nil(t, err)

	dl, err := createRandomDLInDB()
	require.Nil(t, err)

	initialItems := []voucher.VoucherItemInsertRequest{
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

	newItems := []voucher.VoucherItemInsertRequest{
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
	items := []voucher.VoucherItemInsertRequest{
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
	slWithDL, err := createRandomSLInDB(true)
	require.Nil(t, err)

	items := []voucher.VoucherItemInsertRequest{
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
	slWithoutDL, err := createRandomSLInDB(false)
	require.Nil(t, err)

	dl, err := createRandomDLInDB()
	require.Nil(t, err)

	items := []voucher.VoucherItemInsertRequest{
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
	slWithDL, err := createRandomSLInDB(true)
	require.Nil(t, err)

	invalidDLID := generateRandomInt64()

	items := []voucher.VoucherItemInsertRequest{
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

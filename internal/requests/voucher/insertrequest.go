package voucher

type InsertRequest struct {
	Number       string
	VoucherItems []VoucherItemInsertDetail
}

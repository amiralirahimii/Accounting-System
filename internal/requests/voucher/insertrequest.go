package voucher

type VoucherItemInsertRequest struct {
	SLID   int
	DLID   *int
	Debit  int
	Credit int
}

type InsertRequest struct {
	Number       string
	VoucherItems []VoucherItemInsertRequest
}

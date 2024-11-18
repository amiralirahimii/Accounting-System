package voucher

type VoucherItemUpdateDetail struct {
	ID     int
	SLID   int
	DLID   *int
	Debit  int
	Credit int
}

type VoucherItemsUpdate struct {
	Inserted []VoucherItemInsertDetail
	Updated  []VoucherItemUpdateDetail
	Deleted  []int
}

type UpdateRequest struct {
	ID      int
	Number  string
	Version int
	Items   VoucherItemsUpdate
}

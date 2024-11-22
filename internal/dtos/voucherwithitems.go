package dtos

type VoucherWithItemsDto struct {
	ID           int
	Number       string
	RowVersion   int
	VoucherItems []VoucherItemDto
}

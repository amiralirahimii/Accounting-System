package mappers

import (
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/models"
)

func ToVoucherDto(voucher *models.Voucher, voucherItems []models.VoucherItem) *dtos.VoucherDto {
	voucherItemDtos := make([]dtos.VoucherItemDto, len(voucherItems))
	for i, item := range voucherItems {
		voucherItemDtos[i] = dtos.VoucherItemDto{
			ID:     item.ID,
			SLID:   item.SLID,
			DLID:   int(item.DLID.Int64),
			Debit:  item.Debit,
			Credit: item.Credit,
		}
	}

	return &dtos.VoucherDto{
		ID:           voucher.ID,
		Number:       voucher.Number,
		RowVersion:   voucher.RowVersion,
		VoucherItems: voucherItemDtos,
	}
}

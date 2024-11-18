package services

import (
	"accountingsystem/db"
	"accountingsystem/internal/constants"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/voucher"
	"database/sql"
	"errors"
	"log"

	"gorm.io/gorm"
)

type VoucherService struct{}

func (s *VoucherService) CreateVoucher(req *voucher.InsertRequest) (*models.Voucher, error) {
	if err := s.validateVoucherRequest(req); err != nil {
		return nil, err
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return nil, constants.ErrUnexpectedError
	}

	if err := s.checkExistingVoucher(tx, req.Number); err != nil {
		tx.Rollback()
		return nil, err
	}

	voucher, err := s.insertVoucher(tx, req.Number)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.insertVoucherItems(tx, voucher.ID, req.VoucherItems); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return voucher, nil
}

func (s *VoucherService) validateVoucherRequest(req *voucher.InsertRequest) error {
	if req.Number == "" || len(req.Number) > 64 {
		return constants.ErrNumberEmptyOrTooLong
	}
	if len(req.VoucherItems) < 2 || len(req.VoucherItems) > 500 {
		return constants.ErrVoucherItemsCountOutOfRange
	}
	if err := s.validateVoucherItems(req.VoucherItems); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) validateVoucherItems(items []voucher.VoucherItemInsertDetail) error {
	totalDebit := 0
	totalCredit := 0

	for _, item := range items {
		if err := s.validateDebitCredit(item); err != nil {
			return err
		}

		totalDebit += item.Debit
		totalCredit += item.Credit
	}

	if totalDebit != totalCredit {
		return constants.ErrDebitCreditMismatch
	}

	return nil
}

func (s *VoucherService) validateDebitCredit(item voucher.VoucherItemInsertDetail) error {
	isValidDebitCredit := (item.Debit == 0 && item.Credit > 0) || (item.Debit > 0 && item.Credit == 0)
	if !isValidDebitCredit {
		return constants.ErrDebitOrCreditInvalid
	}
	return nil
}

func (s *VoucherService) checkExistingVoucher(tx *gorm.DB, number string) error {
	var existingVoucher models.Voucher
	if err := tx.Where("number = ?", number).First(&existingVoucher).Error; err == nil {
		return constants.ErrVoucherNumberExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Unexpected error while checking voucher number: %v", err)
		return constants.ErrUnexpectedError
	}
	return nil
}

func (s *VoucherService) insertVoucher(tx *gorm.DB, number string) (*models.Voucher, error) {
	voucher := &models.Voucher{
		Number:     number,
		RowVersion: 0,
	}
	if err := tx.Create(voucher).Error; err != nil {
		log.Printf("Error creating voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}
	return voucher, nil
}

func (s *VoucherService) insertVoucherItems(tx *gorm.DB, voucherID int, items []voucher.VoucherItemInsertDetail) error {
	for _, item := range items {
		if err := s.validateSLAndDL(tx, item); err != nil {
			return err
		}
		if err := s.createVoucherItem(tx, voucherID, item); err != nil {
			return err
		}
	}
	return nil
}

func (s *VoucherService) createVoucherItem(tx *gorm.DB, voucherID int, item voucher.VoucherItemInsertDetail) error {
	dlID := s.convertToNullInt64(item.DLID)

	voucherItem := models.VoucherItem{
		VoucherID: voucherID,
		SLID:      item.SLID,
		DLID:      dlID,
		Debit:     item.Debit,
		Credit:    item.Credit,
	}

	if err := tx.Create(&voucherItem).Error; err != nil {
		log.Printf("Error creating voucher item: %v", err)
		return err
	}
	return nil
}

func (s *VoucherService) convertToNullInt64(num *int) sql.NullInt64 {
	if num == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*num), Valid: true}
}

func (s *VoucherService) validateSLAndDL(tx *gorm.DB, item voucher.VoucherItemInsertDetail) error {
	var sl models.SL
	if err := tx.First(&sl, item.SLID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrSLNotFound
		}
		log.Printf("Unexpected error while finding SL: %v", err)
		return constants.ErrUnexpectedError
	}

	if sl.HasDL && item.DLID == nil {
		return constants.ErrDLIDRequired
	} else if !sl.HasDL && item.DLID != nil {
		return constants.ErrDLNotAllowed
	}

	if item.DLID != nil {
		var dl models.DL
		if err := tx.First(&dl, *item.DLID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return constants.ErrDLNotFound
			}
			log.Printf("Unexpected error while finding DL: %v", err)
			return constants.ErrUnexpectedError
		}
	}

	return nil
}

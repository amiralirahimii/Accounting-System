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
	if err := s.validateInsertVoucherRequest(req); err != nil {
		return nil, err
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return nil, constants.ErrUnexpectedError
	}

	voucher, err := s.insertVoucher(tx, req.Number)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if err := s.insertVoucherItems(tx, voucher.ID, req.VoucherItems); err != nil {
		tx.Rollback()
		log.Printf("Error inserting voucher items: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return voucher, nil
}

func (s *VoucherService) validateInsertVoucherRequest(req *voucher.InsertRequest) error {
	if req.Number == "" || len(req.Number) > 64 {
		return constants.ErrNumberEmptyOrTooLong
	}
	if len(req.VoucherItems) < 2 || len(req.VoucherItems) > 500 {
		return constants.ErrVoucherItemsCountOutOfRange
	}
	if err := s.validateVoucherExists(db.DB, req.Number); err != nil {
		return err
	}
	if err := s.validateVoucherItemInsertCreditBalance(req.VoucherItems); err != nil {
		return err
	}
	if err := s.validateVoucherItemInsertDetail(req.VoucherItems); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) validateVoucherItemInsertCreditBalance(items []voucher.VoucherItemInsertDetail) error {
	totalDebit := 0
	totalCredit := 0

	for _, item := range items {
		totalDebit += item.Debit
		totalCredit += item.Credit
	}

	if totalDebit != totalCredit {
		return constants.ErrDebitCreditMismatch
	}

	return nil
}

func (s *VoucherService) validateVoucherItemInsertDetail(items []voucher.VoucherItemInsertDetail) error {
	for _, item := range items {
		if err := s.validateVoucherItem(item.Debit, item.Credit, item.SLID, item.DLID); err != nil {
			return err
		}
	}
	return nil
}

func (s *VoucherService) validateVoucherItem(debit int, credit int, SLID int, DLID *int) error {
	if err := s.validateDebitCredit(debit, credit); err != nil {
		return err
	}
	if err := s.validateSLAndDL(SLID, DLID); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) validateDebitCredit(debit int, credit int) error {
	isValidDebitCredit := (debit == 0 && credit > 0) || (debit > 0 && credit == 0)
	if !isValidDebitCredit {
		return constants.ErrDebitOrCreditInvalid
	}
	return nil
}

func (s *VoucherService) validateVoucherExists(db *gorm.DB, number string) error {
	if _, err := s.getVoucher(db, number); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) getVoucher(db *gorm.DB, number string) (*models.Voucher, error) {
	var existingVoucher models.Voucher
	if err := db.Where("number = ?", number).First(&existingVoucher).Error; err == nil {
		return nil, constants.ErrVoucherNumberExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Unexpected error while checking voucher number: %v", err)
		return nil, constants.ErrUnexpectedError
	}
	return &existingVoucher, nil
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
	var voucherItems []models.VoucherItem

	for _, item := range items {
		dlID := s.convertToNullInt64(item.DLID)
		voucherItems = append(voucherItems, models.VoucherItem{
			VoucherID: voucherID,
			SLID:      item.SLID,
			DLID:      dlID,
			Debit:     item.Debit,
			Credit:    item.Credit,
		})
	}

	if len(voucherItems) > 0 {
		if err := tx.Create(&voucherItems).Error; err != nil {
			log.Printf("Error batch inserting voucher items: %v", err)
			return constants.ErrUnexpectedError
		}
	}

	return nil
}

func (s *VoucherService) convertToNullInt64(num *int) sql.NullInt64 {
	if num == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*num), Valid: true}
}

func (s *VoucherService) validateSLAndDL(SLID int, DLID *int) error {
	var sl models.SL
	if err := db.DB.First(&sl, SLID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrSLNotFound
		}
		log.Printf("Unexpected error while finding SL: %v", err)
		return constants.ErrUnexpectedError
	}

	if sl.HasDL && DLID == nil {
		return constants.ErrDLIDRequired
	} else if !sl.HasDL && DLID != nil {
		return constants.ErrDLNotAllowed
	}

	if DLID != nil {
		var dl models.DL
		if err := db.DB.First(&dl, *DLID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return constants.ErrDLNotFound
			}
			log.Printf("Unexpected error while finding DL: %v", err)
			return constants.ErrUnexpectedError
		}
	}

	return nil
}

func (s *VoucherService) UpdateVoucher(req *voucher.UpdateRequest) (*models.Voucher, error) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return nil, constants.ErrUnexpectedError
	}

	var existingVoucher models.Voucher
	if err := tx.Where("id = ?", req.ID).First(&existingVoucher).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrVoucherNotFound
		}
		log.Printf("Unexpected error while finding voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if existingVoucher.RowVersion != req.Version {
		tx.Rollback()
		return nil, constants.ErrVersionOutdated
	}

	if err := s.validateUpdateVoucherRequest(tx, req); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.applyVoucherChanges(tx, req, &existingVoucher); err != nil {
		tx.Rollback()
		return nil, err
	}

	existingVoucher.Number = req.Number
	existingVoucher.RowVersion++
	if err := tx.Save(&existingVoucher).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return &existingVoucher, nil
}

func (s *VoucherService) applyVoucherChanges(tx *gorm.DB, req *voucher.UpdateRequest, existingVoucher *models.Voucher) error {
	if err := s.deleteVoucherItems(tx, req.Items.Deleted); err != nil {
		return err
	}
	if err := s.insertVoucherItems(tx, existingVoucher.ID, req.Items.Inserted); err != nil {
		return err
	}
	if err := s.updateVoucherItems(tx, req.Items.Updated); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) deleteVoucherItems(tx *gorm.DB, itemIDs []int) error {
	if len(itemIDs) == 0 {
		return nil
	}

	if err := tx.Where("id IN ?", itemIDs).Delete(&models.VoucherItem{}).Error; err != nil {
		log.Printf("Error deleting voucher items: %v", err)
		return constants.ErrUnexpectedError
	}
	return nil
}

func (s *VoucherService) updateVoucherItems(tx *gorm.DB, items []voucher.VoucherItemUpdateDetail) error {
	for _, item := range items {
		if err := s.updateVoucherItem(tx, item); err != nil {
			return err
		}
	}
	return nil
}

func (s *VoucherService) updateVoucherItem(tx *gorm.DB, item voucher.VoucherItemUpdateDetail) error {
	var currentItem models.VoucherItem
	tx.First(&currentItem, item.ID)

	currentItem.SLID = item.SLID
	currentItem.DLID = s.convertToNullInt64(item.DLID)
	currentItem.Debit = item.Debit
	currentItem.Credit = item.Credit

	if err := tx.Save(&currentItem).Error; err != nil {
		log.Printf("Error updating voucher item: %v", err)
		return constants.ErrUnexpectedError
	}
	return nil
}

func (s *VoucherService) validateUpdateVoucherRequest(tx *gorm.DB, req *voucher.UpdateRequest) error {
	if req.Number == "" || len(req.Number) > 64 {
		return constants.ErrNumberEmptyOrTooLong
	}
	if err := s.validateUpdateVoucherVoucherItems(tx, req.Items); err != nil {
		return err
	}
	if err := s.validateVoucherItemsCount(tx, req.Items, req.ID); err != nil {
		return nil
	}
	return nil
}

func (s *VoucherService) validateUpdateVoucherVoucherItems(tx *gorm.DB, items voucher.VoucherItemsUpdate) error {
	for _, item := range items.Inserted {
		if err := s.validateDebitCredit(item.Debit, item.Credit); err != nil {
			return err
		}
		if err := s.validateSLAndDL(item.SLID, item.DLID); err != nil {
			return err
		}
	}
	for _, item := range items.Updated {
		if err := s.validateDebitCredit(item.Debit, item.Credit); err != nil {
			return err
		}
		if err := s.validateSLAndDL(item.SLID, item.DLID); err != nil {
			return err
		}
	}

	if err := s.validateVoucherUpdateDebitCreditBalance(tx, items); err != nil {
		return err
	}

	return nil
}

func (s *VoucherService) validateVoucherUpdateDebitCreditBalance(tx *gorm.DB, items voucher.VoucherItemsUpdate) error {
	totalDebit := 0
	totalCredit := 0

	for _, item := range items.Inserted {
		totalDebit += item.Debit
		totalCredit += item.Credit
	}
	for _, item := range items.Updated {
		var currentItem models.VoucherItem
		if err := tx.First(&currentItem, item.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return constants.ErrVoucherItemNotFound
			}
			log.Printf("Unexpected error while finding voucher item: %v", err)
			return constants.ErrUnexpectedError
		}

		totalDebit += item.Debit - currentItem.Debit
		totalCredit += item.Credit - currentItem.Credit
	}
	for _, itemID := range items.Deleted {
		var currentItem models.VoucherItem
		if err := tx.First(&currentItem, itemID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return constants.ErrVoucherItemNotFound
			}
			log.Printf("Unexpected error while finding voucher item: %v", err)
			return constants.ErrUnexpectedError
		}

		totalDebit -= currentItem.Debit
		totalCredit -= currentItem.Credit
	}

	return nil
}

func (s *VoucherService) validateVoucherItemsCount(tx *gorm.DB, items voucher.VoucherItemsUpdate, voucherID int) error {
	newItemsCount := len(items.Inserted) - len(items.Deleted)
	var existingItemsCount int64 = 0
	if err := tx.Model(&models.VoucherItem{}).Where("voucher_id = ?", voucherID).Count(&existingItemsCount).Error; err != nil {
		log.Printf("Unexpected error while counting voucher items: %v", err)
		return constants.ErrUnexpectedError
	}

	if int(existingItemsCount)+newItemsCount < 2 || int(existingItemsCount)+newItemsCount > 500 {
		return constants.ErrVoucherItemsCountOutOfRange
	}

	return nil
}

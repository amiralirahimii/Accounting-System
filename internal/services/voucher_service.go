package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/mappers"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/voucher"
	"database/sql"
	"errors"
	"log"

	"gorm.io/gorm"
)

type VoucherService struct {
	db *gorm.DB
}

func (s *VoucherService) CreateVoucher(req *voucher.InsertRequest) (*dtos.VoucherWithItemsDto, error) {
	if err := s.validateInsertVoucherRequest(req); err != nil {
		return nil, err
	}

	tx := s.db.Begin()
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

	voucherItems, err := s.insertVoucherItems(tx, voucher.ID, req.VoucherItems)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting voucher items: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return mappers.ToVoucherWithItemsDto(voucher, voucherItems), nil
}

func (s *VoucherService) validateInsertVoucherRequest(req *voucher.InsertRequest) error {
	if err := s.validateNumber(req.Number); err != nil {
		return err
	}
	if err := s.validateVoucherItemsCountInInsertRequest(req.VoucherItems); err != nil {
		return err
	}
	if err := s.validateVoucherNumberIsUnique(req.Number); err != nil {
		return err
	}
	if err := s.validateVoucherItemInsertCreditBalance(req.VoucherItems); err != nil {
		return err
	}
	if err := s.validateVoucherItemInsertDetails(req.VoucherItems); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) validateNumber(number string) error {
	if number == "" || len(number) > 64 {
		return constants.ErrNumberEmptyOrTooLong
	}
	return nil
}

func (s *VoucherService) validateVoucherItemsCountInInsertRequest(items []voucher.VoucherItemInsertDetail) error {
	if err := s.validateCorrectVoucherItemNumber(len(items)); err != nil {
		return err
	}
	return nil
}

func (s *VoucherService) validateCorrectVoucherItemNumber(number int) error {
	if number < 2 || number > 500 {
		return constants.ErrVoucherItemsCountOutOfRange
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

func (s *VoucherService) validateVoucherItemInsertDetails(items []voucher.VoucherItemInsertDetail) error {
	for _, item := range items {
		if err := s.validateDebitCredit(item.Debit, item.Credit); err != nil {
			return err
		}
		if err := s.validateSLAndDL(item.SLID, item.DLID); err != nil {
			return err
		}
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

func (s *VoucherService) validateVoucherNumberIsUnique(number string) error {
	var existingVoucher models.Voucher
	if err := s.db.Where("number = ?", number).First(&existingVoucher).Error; err == nil {
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

func (s *VoucherService) insertVoucherItems(tx *gorm.DB, voucherID int, items []voucher.VoucherItemInsertDetail) ([]models.VoucherItem, error) {
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
			return nil, constants.ErrUnexpectedError
		}
	}

	return voucherItems, nil
}

func (s *VoucherService) convertToNullInt64(num *int) sql.NullInt64 {
	if num == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*num), Valid: true}
}

func (s *VoucherService) validateSLAndDL(SLID int, DLID *int) error {
	sl, err := s.validateSLExists(SLID)
	if err != nil {
		return err
	}

	if err := s.validateDLRequirement(sl.HasDL, DLID); err != nil {
		return err
	}

	if DLID != nil {
		if err := s.validateDLExists(*DLID); err != nil {
			return err
		}
	}

	return nil
}

func (s *VoucherService) validateSLExists(SLID int) (*models.SL, error) {
	var sl models.SL
	if err := s.db.First(&sl, SLID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrSLNotFound
		}
		log.Printf("Unexpected error while finding SL: %v", err)
		return nil, constants.ErrUnexpectedError
	}
	return &sl, nil
}

func (s *VoucherService) validateDLExists(DLID int) error {
	var dl models.DL
	if err := s.db.First(&dl, DLID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrDLNotFound
		}
		log.Printf("Unexpected error while finding DL: %v", err)
		return constants.ErrUnexpectedError
	}
	return nil
}

func (s *VoucherService) validateDLRequirement(SLhasDL bool, DLID *int) error {
	if SLhasDL && DLID == nil {
		return constants.ErrDLIDRequired
	} else if !SLhasDL && DLID != nil {
		return constants.ErrDLNotAllowed
	}
	return nil
}

func (s *VoucherService) UpdateVoucher(req *voucher.UpdateRequest) (*dtos.VoucherDto, error) {
	targetVoucher, err := s.validateUpdateVoucherRequest(req)
	if err != nil {
		return nil, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return nil, constants.ErrUnexpectedError
	}

	if err := s.applyVoucherItemChanges(tx, req, targetVoucher); err != nil {
		tx.Rollback()
		return nil, err
	}

	targetVoucher.Number = req.Number
	targetVoucher.RowVersion++

	if err := tx.Save(targetVoucher).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return mappers.ToVoucherDto(targetVoucher), nil
}

func (s *VoucherService) validateUpdateVoucherRequest(req *voucher.UpdateRequest) (*models.Voucher, error) {
	if err := s.validateNumber(req.Number); err != nil {
		return nil, err
	}
	targetVoucher, err := s.validateVoucherExists(req.ID)
	if err != nil {
		return nil, err
	}
	if err := s.validateVersion(targetVoucher.RowVersion, req.Version); err != nil {
		return nil, err
	}
	if err := s.validateVoucherItemsCountInUpdateRequest(req.Items, req.ID); err != nil {
		return nil, err
	}
	if err := s.validateVoucherItemsInUpdateRequest(req.Items); err != nil {
		return nil, err
	}
	return targetVoucher, nil
}

func (s *VoucherService) validateVersion(existingVersion int, requestVersion int) error {
	if existingVersion != requestVersion {
		return constants.ErrVersionOutdated
	}
	return nil
}

func (s *VoucherService) validateVoucherExists(id int) (*models.Voucher, error) {
	var existingVoucher models.Voucher
	if err := s.db.Where("id = ?", id).First(&existingVoucher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrVoucherNotFound
		}
		log.Printf("Unexpected error while finding voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}
	return &existingVoucher, nil
}

func (s *VoucherService) validateVoucherItemsInUpdateRequest(items voucher.VoucherItemsUpdate) error {
	if err := s.validateVoucherItemInsertDetails(items.Inserted); err != nil {
		return err
	}
	if err := s.validateVoucherItemUpdateDetails(items.Updated); err != nil {
		return err
	}
	if err := s.validateVoucherItemDeleteDetails(items.Deleted); err != nil {
		return err
	}
	if err := s.validateVoucherUpdateDebitCreditBalance(items); err != nil {
		return err
	}

	return nil
}

func (s *VoucherService) validateVoucherItemUpdateDetails(items []voucher.VoucherItemUpdateDetail) error {
	for _, item := range items {
		if err := s.validateVoucherItemExists(item.ID); err != nil {
			return err
		}
		if err := s.validateDebitCredit(item.Debit, item.Credit); err != nil {
			return err
		}
		if err := s.validateSLAndDL(item.SLID, item.DLID); err != nil {
			return err
		}
	}
	return nil
}

func (s *VoucherService) validateVoucherItemDeleteDetails(items []int) error {
	for _, itemID := range items {
		if err := s.validateVoucherItemExists(itemID); err != nil {
			return err
		}
	}
	return nil
}

func (s *VoucherService) validateVoucherItemExists(itemID int) error {
	var existingItem models.VoucherItem
	if err := s.db.First(&existingItem, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrVoucherItemNotFound
		}
		log.Printf("Unexpected error while finding voucher item: %v", err)
		return constants.ErrUnexpectedError
	}
	return nil
}

func (s *VoucherService) validateVoucherUpdateDebitCreditBalance(items voucher.VoucherItemsUpdate) error {
	totalDebitAddedInInsert, totalCreditAddedInInsert := s.calculateInsertedBalances(items)
	totalDebitAddedInUpdate, totalCreditAddedInUpdate := s.calculateUpdatedBalances(items)
	totalDebitAddedInDelete, totalCreditAddedInDelete := s.calculateDeletedBalances(items)

	totalDebitAdded := totalDebitAddedInInsert + totalDebitAddedInUpdate + totalDebitAddedInDelete
	totalCreditAdded := totalCreditAddedInInsert + totalCreditAddedInUpdate + totalCreditAddedInDelete

	if totalDebitAdded != totalCreditAdded {
		return constants.ErrDebitCreditMismatch
	}

	return nil
}

func (s *VoucherService) calculateInsertedBalances(items voucher.VoucherItemsUpdate) (int, int) {
	totalDebit := 0
	totalCredit := 0
	for _, item := range items.Inserted {
		totalDebit += item.Debit
		totalCredit += item.Credit
	}
	return totalDebit, totalCredit
}

func (s *VoucherService) calculateUpdatedBalances(items voucher.VoucherItemsUpdate) (int, int) {
	totalDebit := 0
	totalCredit := 0
	for _, item := range items.Updated {
		var currentItem models.VoucherItem
		s.db.First(&currentItem, item.ID)
		totalDebit += item.Debit - currentItem.Debit
		totalCredit += item.Credit - currentItem.Credit
	}
	return totalDebit, totalCredit
}

func (s *VoucherService) calculateDeletedBalances(items voucher.VoucherItemsUpdate) (int, int) {
	totalDebit := 0
	totalCredit := 0
	for _, itemID := range items.Deleted {
		var currentItem models.VoucherItem
		s.db.First(&currentItem, itemID)
		totalDebit -= currentItem.Debit
		totalCredit -= currentItem.Credit
	}
	return totalDebit, totalCredit
}

func (s *VoucherService) validateVoucherItemsCountInUpdateRequest(items voucher.VoucherItemsUpdate, voucherID int) error {
	newItemsCount := len(items.Inserted) - len(items.Deleted)
	var existingItemsCount int64 = 0
	if err := s.db.Model(&models.VoucherItem{}).Where("voucher_id = ?", voucherID).Count(&existingItemsCount).Error; err != nil {
		log.Printf("Unexpected error while counting voucher items: %v", err)
		return constants.ErrUnexpectedError
	}

	if err := s.validateCorrectVoucherItemNumber(int(existingItemsCount) + newItemsCount); err != nil {
		return err
	}

	return nil
}

func (s *VoucherService) applyVoucherItemChanges(tx *gorm.DB, req *voucher.UpdateRequest, existingVoucher *models.Voucher) error {
	if err := s.deleteVoucherItems(tx, req.Items.Deleted); err != nil {
		return err
	}
	if _, err := s.insertVoucherItems(tx, existingVoucher.ID, req.Items.Inserted); err != nil {
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

func (s *VoucherService) DeleteVoucher(req *voucher.DeleteRequest) error {
	targetVoucher, err := s.validateDeleteVoucherRequest(req)
	if err != nil {
		return err
	}

	if err := s.db.Delete(&targetVoucher).Error; err != nil {
		return constants.ErrUnexpectedError
	}

	return nil
}

func (s *VoucherService) validateDeleteVoucherRequest(req *voucher.DeleteRequest) (*models.Voucher, error) {
	targetVoucher, err := s.validateVoucherExists(req.ID)
	if err != nil {
		return nil, err
	}
	if err := s.validateVersion(targetVoucher.RowVersion, req.Version); err != nil {
		return nil, err
	}
	return targetVoucher, nil
}

func (s *VoucherService) GetVoucher(req *voucher.GetRequest) (*dtos.VoucherWithItemsDto, error) {
	targetVoucher, err := s.validateGetVoucherRequest(req)
	if err != nil {
		return nil, err
	}

	var voucherItems []models.VoucherItem
	if err := s.db.Where("voucher_id = ?", targetVoucher.ID).Find(&voucherItems).Error; err != nil {
		log.Printf("unexpected error while finding voucher items: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return mappers.ToVoucherWithItemsDto(targetVoucher, voucherItems), nil
}

func (s *VoucherService) validateGetVoucherRequest(req *voucher.GetRequest) (*models.Voucher, error) {
	targetVoucher, err := s.validateVoucherExists(req.ID)
	if err != nil {
		return nil, err
	}

	return targetVoucher, nil
}

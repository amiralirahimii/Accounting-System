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

func (s *VoucherService) applyVoucherCreation(req *voucher.InsertRequest) (*dtos.VoucherWithItemsDto, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	voucher, err := s.insertVoucher(tx, req.Number)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	voucherItems, err := s.insertVoucherItems(tx, voucher.ID, req.VoucherItems)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return mappers.ToVoucherWithItemsDto(voucher, voucherItems), nil
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
			return nil, err
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

func (s *VoucherService) validateVoucherNumberIsUnique(number string) error {
	var existingVoucher models.Voucher
	if err := s.db.Where("number = ?", number).First(&existingVoucher).Error; err == nil {
		return constants.ErrVoucherNumberExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (s *VoucherService) validateVoucherItemInsertCreditBalance(items []voucher.VoucherItemInsertDetail) error {
	totalDebit, totalCredit := s.calculateInsertedBalances(items)

	if totalDebit != totalCredit {
		return constants.ErrDebitCreditMismatch
	}

	return nil
}

func (s *VoucherService) calculateInsertedBalances(items []voucher.VoucherItemInsertDetail) (int, int) {
	totalDebit := 0
	totalCredit := 0
	for _, item := range items {
		totalDebit += item.Debit
		totalCredit += item.Credit
	}
	return totalDebit, totalCredit
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
		return nil, err
	}
	return &sl, nil
}

func (s *VoucherService) validateDLRequirement(SLhasDL bool, DLID *int) error {
	if SLhasDL && DLID == nil {
		return constants.ErrDLIDRequired
	} else if !SLhasDL && DLID != nil {
		return constants.ErrDLNotAllowed
	}
	return nil
}

func (s *VoucherService) validateDLExists(DLID int) error {
	var dl models.DL
	if err := s.db.First(&dl, DLID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrDLNotFound
		}
		return err
	}
	return nil
}

func (s *VoucherService) applyVoucherUpdate(req *voucher.UpdateRequest, targetVoucher *models.Voucher) (*dtos.VoucherDto, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	targetVoucher.Number = req.Number
	targetVoucher.RowVersion++

	if err := tx.Save(targetVoucher).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.applyVoucherItemChanges(tx, req, targetVoucher); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return mappers.ToVoucherDto(targetVoucher), nil
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
		return err
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
		return err
	}
	return nil
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

func (s *VoucherService) validateVoucherExists(id int) (*models.Voucher, error) {
	var existingVoucher models.Voucher
	if err := s.db.Where("id = ?", id).First(&existingVoucher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrVoucherNotFound
		}
		return nil, err
	}
	return &existingVoucher, nil
}

func (s *VoucherService) validateVersion(existingVersion int, requestVersion int) error {
	if existingVersion != requestVersion {
		return constants.ErrVersionOutdated
	}
	return nil
}

func (s *VoucherService) validateVoucherItemsCountInUpdateRequest(items voucher.VoucherItemsUpdate, voucherID int) error {
	newItemsCount := len(items.Inserted) - len(items.Deleted)
	var existingItemsCount int64 = 0
	if err := s.db.Model(&models.VoucherItem{}).Where("voucher_id = ?", voucherID).Count(&existingItemsCount).Error; err != nil {
		return err
	}

	if err := s.validateCorrectVoucherItemNumber(int(existingItemsCount) + newItemsCount); err != nil {
		return err
	}

	return nil
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

func (s *VoucherService) validateVoucherItemExists(itemID int) error {
	var existingItem models.VoucherItem
	if err := s.db.First(&existingItem, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrVoucherItemNotFound
		}
		return err
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

func (s *VoucherService) validateVoucherUpdateDebitCreditBalance(items voucher.VoucherItemsUpdate) error {
	totalDebitAddedInInsert, totalCreditAddedInInsert := s.calculateInsertedBalances(items.Inserted)
	totalDebitAddedInUpdate, totalCreditAddedInUpdate := s.calculateUpdatedBalances(items.Updated)
	totalDebitAddedInDelete, totalCreditAddedInDelete := s.calculateDeletedBalances(items.Deleted)

	totalDebitAdded := totalDebitAddedInInsert + totalDebitAddedInUpdate + totalDebitAddedInDelete
	totalCreditAdded := totalCreditAddedInInsert + totalCreditAddedInUpdate + totalCreditAddedInDelete

	if totalDebitAdded != totalCreditAdded {
		return constants.ErrDebitCreditMismatch
	}

	return nil
}

func (s *VoucherService) calculateUpdatedBalances(items []voucher.VoucherItemUpdateDetail) (int, int) {
	totalDebit := 0
	totalCredit := 0
	for _, item := range items {
		var currentItem models.VoucherItem
		s.db.First(&currentItem, item.ID)
		totalDebit += item.Debit - currentItem.Debit
		totalCredit += item.Credit - currentItem.Credit
	}
	return totalDebit, totalCredit
}

func (s *VoucherService) calculateDeletedBalances(items []int) (int, int) {
	totalDebit := 0
	totalCredit := 0
	for _, itemID := range items {
		var currentItem models.VoucherItem
		s.db.First(&currentItem, itemID)
		totalDebit -= currentItem.Debit
		totalCredit -= currentItem.Credit
	}
	return totalDebit, totalCredit
}

func (s *VoucherService) applyVoucherDeletion(targetVoucher *models.Voucher) error {
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

func (s *VoucherService) applyVoucherGet(targetVoucher *models.Voucher) (*dtos.VoucherWithItemsDto, error) {
	var voucherItems []models.VoucherItem
	if err := s.db.Where("voucher_id = ?", targetVoucher.ID).Find(&voucherItems).Error; err != nil {
		return nil, err
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

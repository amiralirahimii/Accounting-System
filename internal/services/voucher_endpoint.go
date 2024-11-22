package services

import (
	"accountingsystem/internal/constants"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/requests/voucher"
	"log"

	"gorm.io/gorm"
)

type VoucherService struct {
	db *gorm.DB
}

func (s *VoucherService) InitService(db *gorm.DB) {
	s.db = db
}

func (s *VoucherService) CreateVoucher(req *voucher.InsertRequest) (*dtos.VoucherWithItemsDto, error) {
	if err := s.validateInsertVoucherRequest(req); err != nil {
		return nil, err
	}

	voucherWithItemsDto, err := s.applyVoucherCreation(req)
	if err != nil {
		log.Printf("unexpected error while creating voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return voucherWithItemsDto, nil
}

func (s *VoucherService) UpdateVoucher(req *voucher.UpdateRequest) (*dtos.VoucherDto, error) {
	targetVoucher, err := s.validateUpdateVoucherRequest(req)
	if err != nil {
		return nil, err
	}

	voucherDto, err := s.applyVoucherUpdate(req, targetVoucher)
	if err != nil {
		log.Printf("unexpected error while updating voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return voucherDto, nil
}

func (s *VoucherService) DeleteVoucher(req *voucher.DeleteRequest) error {
	targetVoucher, err := s.validateDeleteVoucherRequest(req)
	if err != nil {
		return err
	}

	if err := s.applyVoucherDeletion(targetVoucher); err != nil {
		log.Printf("unexpected error while deleting voucher: %v", err)
		return constants.ErrUnexpectedError
	}

	return nil
}

func (s *VoucherService) GetVoucher(req *voucher.GetRequest) (*dtos.VoucherWithItemsDto, error) {
	targetVoucher, err := s.validateGetVoucherRequest(req)
	if err != nil {
		return nil, err
	}

	voucherWithItemsDto, err := s.applyVoucherGet(targetVoucher)
	if err != nil {
		log.Printf("unexpected error while getting voucher: %v", err)
		return nil, constants.ErrUnexpectedError
	}

	return voucherWithItemsDto, nil
}

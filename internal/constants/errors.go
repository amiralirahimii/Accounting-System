package constants

import "errors"

var (
	ErrUnexpectedError             = errors.New("Something went wrong")
	ErrEnvNotFound                 = errors.New("environment variable not found")
	ErrCodeEmptyOrTooLong          = errors.New("code cannot be empty or more than 64 characters")
	ErrTitleEmptyOrTooLong         = errors.New("title cannot be empty or more than 64 characters")
	ErrCodeAlreadyExists           = errors.New("code should be unique")
	ErrTitleAlreadyExists          = errors.New("title should be unique")
	ErrDLNotFound                  = errors.New("DL not found")
	ErrVersionOutdated             = errors.New("version is outdated")
	ErrSLNotFound                  = errors.New("SL not found")
	ErrNumberEmptyOrTooLong        = errors.New("number cannot be empty or more than 64 characters")
	ErrVoucherNumberExists         = errors.New("voucher number already exists")
	ErrVoucherItemsCountOutOfRange = errors.New("voucher items count should be between 2 and 500")
	ErrDebitOrCreditInvalid        = errors.New("one and only one of debit or credit should be greater than 0")
	ErrDLIDRequired                = errors.New("Provided SL requires DL")
	ErrDLNotAllowed                = errors.New("Provided SL does not require DL")
	ErrDebitCreditMismatch         = errors.New("debits ad credits should be equal in a voucher")
)

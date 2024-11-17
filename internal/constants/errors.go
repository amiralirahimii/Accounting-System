package constants

import "errors"

var (
	ErrUnexpectedError     = errors.New("Something went wrong")
	ErrEnvNotFound         = errors.New("environment variable not found")
	ErrCodeEmptyOrTooLong  = errors.New("code cannot be empty or more than 64 characters")
	ErrTitleEmptyOrTooLong = errors.New("title cannot be empty or more than 64 characters")
	ErrCodeAlreadyExists   = errors.New("code should be unique")
	ErrTitleAlreadyExists  = errors.New("title should be unique")
	ErrDLNotFound          = errors.New("DL not found")
	ErrVersionOutdated     = errors.New("version is outdated")
	ErrSLNotFound          = errors.New("SL not found")
)

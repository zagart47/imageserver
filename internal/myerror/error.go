package myerror

import "errors"

type Error struct {
	NotExists    error
	UpdateFailed error
	FileNotFound error
	MdError      error
	BuffError    error
}

func NewErrors() Error {
	return Error{
		NotExists:    errors.New("row not exists"),
		UpdateFailed: errors.New("update failed"),
		FileNotFound: errors.New("file not found"),
		MdError:      errors.New("metadata incoming error"),
		BuffError:    errors.New("buffer reading error"),
	}
}

var Err = NewErrors()

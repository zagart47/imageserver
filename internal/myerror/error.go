package myerror

import "errors"

type Error struct {
	NotExists    error
	UpdateFailed error
	FileNotFound error
	MetaData     error
	Buffer       error
}

func NewErrors() Error {
	return Error{
		NotExists:    errors.New("row not exists"),
		UpdateFailed: errors.New("update failed"),
		FileNotFound: errors.New("file not found"),
		MetaData:     errors.New("metadata incoming error"),
		Buffer:       errors.New("buffer reading error"),
	}
}

var Err = NewErrors()

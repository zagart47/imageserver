package myerror

import "errors"

var Err = NewErrors()

type Error struct {
	NotExists    error
	UpdateFailed error
	FileNotFound error
	Metadata     error
	Buffer       error
	InvFileName  error
}

func NewErrors() Error {
	return Error{
		NotExists:    errors.New("row not exists"),
		UpdateFailed: errors.New("update failed"),
		FileNotFound: errors.New("file not found"),
		Metadata:     errors.New("metadata incoming error"),
		Buffer:       errors.New("buffer reading error"),
		InvFileName:  errors.New("invalid filename"),
	}
}

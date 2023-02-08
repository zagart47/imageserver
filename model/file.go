package model

import (
	"bytes"
	"fmt"
)

type File struct {
	Name   string
	Path   string
	Buffer *bytes.Buffer
}

func NewFile(name string) *File {
	dir := "files"
	pathToSave := fmt.Sprintf("%s/%s", dir, name)
	return &File{
		Name:   name,
		Path:   pathToSave,
		Buffer: &bytes.Buffer{},
	}
}

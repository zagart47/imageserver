package storage

import (
	"bytes"
	"strings"
)

type File struct {
	Name   string
	Path   string
	Buffer *bytes.Buffer
}

func NewFile(name string) File {
	dir := "files"
	els := []string{dir, name}
	pathToSave := strings.Join(els, "/")
	return File{
		Name:   name,
		Path:   pathToSave,
		Buffer: &bytes.Buffer{},
	}
}

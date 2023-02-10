package file

type File struct {
	FileName string
	Created  string
	Updated  string
}

type ListFile []File

func NewListFile(name string) ListFile {
	return ListFile{File{
		FileName: name,
		Created:  "",
		Updated:  "",
	}}
}

func NewFile(name string) *File {
	return &File{
		FileName: name,
		Created:  "",
		Updated:  "",
	}
}

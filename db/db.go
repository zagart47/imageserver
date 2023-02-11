package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"imageserver/file"
	"strings"
	"time"
)

var DB, _ = sql.Open("sqlite3", "db/files.db")

var (
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrFileFound    = errors.New("file found")
	ErrFileNotFound = errors.New("file not found")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() *SQLiteRepository {
	return &SQLiteRepository{
		db: DB,
	}
}

var Repo = NewSQLiteRepository()

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS files(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        file_name TEXT NOT NULL,
        created TEXT NOT NULL,
        updated TEXT NOT NULL
    );
    `
	_, err := r.db.Exec(query)
	return err
}
func (r *SQLiteRepository) Create(filename string) error {
	_, err := r.db.Exec("INSERT INTO files(file_name, created, updated) values(?,?,?)", filename, CurrentTime(), "have no update")
	if err != nil {
		return ErrUpdateFailed
	}
	return nil
}

func (r *SQLiteRepository) Update(filename string) error {
	_, err := r.db.Exec("UPDATE files SET updated = ? WHERE file_name = ?", CurrentTime(), filename)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) CheckFileName(filename string) error {
	if filename == "" {
		return errors.New("invalid updated filename")
	}
	err := r.Migrate()
	if err != nil {
		return err
	}
	all, err := r.All()
	if err != nil {
		return err
	}
	for _, v := range all {
		if v.FileName == filename {
			return ErrFileFound
		}
	}
	return ErrFileNotFound
}

func (r *SQLiteRepository) All() (file.ListFile, error) {
	err := r.Migrate()
	if err != nil {
		return file.ListFile{}, err
	}
	rows, err := r.db.Query("SELECT file_name, created, updated FROM files")
	if err != nil {
		return file.ListFile{}, ErrNotExists
	}
	defer rows.Close()

	var all file.ListFile
	for rows.Next() {
		var file file.File
		err = rows.Scan(&file.FileName, &file.Created, &file.Updated)
		if err != nil {
			return all, err
		}
		all = append(all, file)
	}
	return all, nil
}

func CurrentTime() string {
	return time.Now().Format("02.01.2006 15:04:05")
}

func DownloadFileList() string {
	repo := NewSQLiteRepository()
	all, err := repo.All()
	if err != nil {
		return ""
	}
	fl := file.ListFile{}
	for _, v := range all {
		fl = append(fl, file.File{
			FileName: v.FileName,
			Created:  v.Created,
			Updated:  v.Updated,
		})
	}
	return MakeTable(fl)
}

type Lengths struct {
	FileNameLength int
	CreatedLength  int
	UpdatedLength  int
}

func NewLengths(lf file.ListFile) *Lengths {
	sl := MaxStringLengths(lf)
	return &Lengths{
		FileNameLength: sl.FileNameLength,
		CreatedLength:  sl.CreatedLength,
		UpdatedLength:  sl.UpdatedLength,
	}
}

func MakeTable(lf file.ListFile) string {
	l := NewLengths(lf)
	upHead := fmt.Sprintf("%c%s%c%s%c%s%c", te.LeftUp, te.RepeatLine(l.FileNameLength), te.MiddleUp, te.RepeatLine(l.CreatedLength), te.MiddleUp, te.RepeatLine(l.UpdatedLength), te.RightUp)
	midHead := fmt.Sprintf("%c%s%c%s%c%s%c", te.V, Fitting(te.tName, l.FileNameLength), te.V, Fitting(te.tCreated, l.CreatedLength), te.V, Fitting(te.tUpdated, l.UpdatedLength), te.V)
	downHead := fmt.Sprintf("%c%s%c%s%c%s%c", te.LeftMiddle, te.RepeatLine(l.FileNameLength), te.CenterMiddle, te.RepeatLine(l.CreatedLength), te.CenterMiddle, te.RepeatLine(l.UpdatedLength), te.RightMiddle)
	var table string
	for _, v := range lf {
		table = fmt.Sprintf("%s%c%s%c%s%c%s%c\n", table, te.V, Fitting(v.FileName, l.FileNameLength), te.V, Fitting(v.Created, l.CreatedLength), te.V, Fitting(v.Updated, l.UpdatedLength), te.V)
	}
	footer := fmt.Sprintf("%c%s%c%s%c%s%c", te.LeftBottom, te.RepeatLine(l.FileNameLength), te.MiddleBottom, te.RepeatLine(l.CreatedLength), te.MiddleBottom, te.RepeatLine(l.UpdatedLength), te.RightBottom)
	result := fmt.Sprintf("%s\n%s\n%s\n%s%s", upHead, midHead, downHead, table, footer)
	return result
}

func MaxStringLengths(lf file.ListFile) Lengths {
	var maxFileNameLength, maxCreatedLength, maxUpdatedLength int
	air := 2
	for _, v := range lf {
		if len(v.FileName) > maxFileNameLength {
			maxFileNameLength = len(v.FileName)
		}
		if len(v.Created) > maxCreatedLength {
			maxCreatedLength = len(v.Created)
		}
		if len(v.Updated) > maxUpdatedLength {
			maxUpdatedLength = len(v.Updated)
		}
	}
	if maxFileNameLength < len(te.tName) {
		maxFileNameLength = len(te.tName)
	}
	if maxCreatedLength < len(te.tCreated) {
		maxCreatedLength = len(te.tCreated)
	}
	if maxUpdatedLength < len(te.tUpdated) {
		maxUpdatedLength = len(te.tUpdated)
	}
	return Lengths{maxFileNameLength + air, maxCreatedLength + air, maxUpdatedLength + air}
}

func Fitting(s string, n int) string {
	for len(s) < n {
		s = fmt.Sprintf("%s%c", s, te.WhiteSpace)
		if len(s) == n {
			break
		}
		s = fmt.Sprintf("%c%s", te.WhiteSpace, s)
		if len(s) == n {
			break
		}
	}
	return s
}

type CP struct {
	WhiteSpace   rune
	LeftUp       rune
	MiddleUp     rune
	RightUp      rune
	V            rune
	H            rune
	LeftMiddle   rune
	CenterMiddle rune
	RightMiddle  rune
	LeftBottom   rune
	MiddleBottom rune
	RightBottom  rune
	tName        string
	tCreated     string
	tUpdated     string
}

var te = NewCP()

func (te CP) RepeatLine(n int) string {
	return strings.Repeat(string(te.H), n)
}

func NewCP() *CP {
	return &CP{
		WhiteSpace:   '\u0020',
		LeftUp:       '\u2554',
		MiddleUp:     '\u2566',
		RightUp:      '\u2557',
		V:            '\u2551',
		H:            '\u2550',
		LeftMiddle:   '\u2560',
		CenterMiddle: '\u256c',
		RightMiddle:  '\u2563',
		LeftBottom:   '\u255a',
		MiddleBottom: '\u2569',
		RightBottom:  '\u255d',
		tName:        "File name",
		tCreated:     "Created",
		tUpdated:     "Last updated",
	}
}

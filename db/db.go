package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"imageserver/file"
	"time"
)

var DB, _ = sql.Open("sqlite3", "db/files.db")

var (
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrFileNotFound = errors.New("file not found")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

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
	if len(filename) == 0 {
		return errors.New("invalid updated filename")
	}
	if err := r.Migrate(); err != nil {
		return err
	}
	all, err := r.AllRecords()
	if err != nil {
		return err
	}
	for _, v := range all {
		if v.FileName == filename {
			return nil
		}
	}
	return ErrFileNotFound
}

func (r *SQLiteRepository) AllRecords() (file.ListFile, error) {
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

func (r *SQLiteRepository) DownloadFileList() (*file.ListFile, error) {
	all, err := r.AllRecords()
	if err != nil {
		return nil, err
	}
	fl := file.ListFile{}
	for _, v := range all {
		fl = append(fl, file.File{
			FileName: v.FileName,
			Created:  v.Created,
			Updated:  v.Updated,
		})
	}
	return &fl, nil
}

package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	uploadpb "imageserver/pkg/proto"
	"time"
)

var DB, _ = sql.Open("sqlite3", "db/files.db")

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() *SQLiteRepository {
	return &SQLiteRepository{
		db: DB,
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
	for _, v := range all {
		if v.FileName == filename {
			return errors.New("file found")
		}
	}
	return errors.New("file not found")
}
func (r *SQLiteRepository) Create(filename string) error {
	_, err := r.db.Exec("INSERT INTO files(file_name, created, updated) values(?,?,?)", filename, CurrentTime(), "have no update")
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) All() ([]*uploadpb.File, error) {
	err := r.Migrate()
	if err != nil {
		return nil, err
	}
	rows, _ := r.db.Query("SELECT file_name, created, updated FROM files")
	defer rows.Close()

	var all []*uploadpb.File
	for rows.Next() {
		var file uploadpb.File
		_ = rows.Scan(&file.FileName, &file.Created, &file.Updated)
		all = append(all, &file)
	}
	return all, nil
}

func CurrentTime() string {
	return time.Now().Format("02.01.2006 15:04:05")
}

package repository

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"imageserver/internal/config"
	"imageserver/internal/model"
	"imageserver/internal/myerror"
	"time"
)

var DB, _ = sql.Open(config.SqlConnect.DriverName, config.SqlConnect.DataSourceName)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) SQLiteRepository {
	return SQLiteRepository{
		db: db,
	}
}

// Migrate prepares the table of database for work.
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

// Create adds the file name with the current date and time to the table of database.
func (r *SQLiteRepository) Create(filename string) error {
	_, err := r.db.Exec("INSERT INTO files(file_name, created, updated) values(?,?,?)", filename, CurrentTime(), "have no update")
	if err != nil {
		return myerror.Err.UpdateFailed
	}
	return nil
}

// Update adds the current date and time to the "updated" column.
func (r *SQLiteRepository) Update(filename string) error {
	_, err := r.db.Exec("UPDATE files SET updated = ? WHERE file_name = ?", CurrentTime(), filename)
	if err != nil {
		return err
	}
	return nil
}

// CheckFileName searches for the passed file name in the database table.
// If file name is found, Update is called.
// If no file name is found, then Create is called.
func (r *SQLiteRepository) CheckFileName(filename string) error {
	if len(filename) == 0 {
		return errors.New("invalid updated filename")
	}
	if err := r.Migrate(); err != nil {
		return err
	}
	all, err := r.ShowAllRecords()
	if err != nil {
		return err
	}
	for _, v := range all {
		if v.FileName == filename {
			return nil
		}
	}
	return myerror.Err.FileNotFound
}

// ShowAllRecords generates a record structure from all available rows in the database table.
func (r *SQLiteRepository) ShowAllRecords() (model.ListFile, error) {
	err := r.Migrate()
	if err != nil {
		return model.ListFile{}, err
	}
	rows, err := r.db.Query("SELECT file_name, created, updated FROM files")
	if err != nil {
		return model.ListFile{}, myerror.Err.NotExists
	}
	defer rows.Close()

	var all model.ListFile
	for rows.Next() {
		var f model.File
		err = rows.Scan(&f.FileName, &f.Created, &f.Updated)
		if err != nil {
			return all, err
		}
		all = append(all, f)
	}
	return all, nil
}

// CurrentTime returns the current date and time.
func CurrentTime() string {
	return time.Now().Format("02.01.2006 15:04:05")
}

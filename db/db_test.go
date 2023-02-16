package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"imageserver/file"
	"os"
	"testing"
	"time"
)

func TestCurrentTime(t *testing.T) {
	var s time.Time
	require.IsType(t, CurrentTime(), s.Format("02.01.2006 15:04:05"))
}

func TestNewSQLiteRepository(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	require.NoError(t, err)
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	var sq *SQLiteRepository
	require.IsType(t, NewSQLiteRepository(db), sq)
}

func TestDB_All(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	require.NoError(t, err)
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	require.NoError(t, err)

	_, err = db.Exec("drop table files")
	r.Migrate()
	err = r.Create("test")
	require.NoError(t, err)

	t.Run("1test", func(t *testing.T) {
		_, err := r.AllRecords()
		require.NoError(t, err)
	})

	t.Run("2test", func(t *testing.T) {
		_, err = db.Exec("drop table files")
		require.NoError(t, err)
		_, err := r.AllRecords()
		require.NoError(t, err)
	})

	db.Close()
	t.Run("3test", func(t *testing.T) {
		_, err := r.AllRecords()
		require.Error(t, err)
	})
}

func TestDB_CheckFileName(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	require.NoError(t, err)
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	require.NoError(t, err)
	_, err = db.Exec("DROP TABLE files")
	err = r.Migrate()
	require.NoError(t, err)
	err = r.Create("test")
	require.NoError(t, err)

	t.Run("1test", func(t *testing.T) {
		err := r.CheckFileName("test")
		require.NoError(t, err)
	})
	t.Run("2test", func(t *testing.T) {
		err := r.CheckFileName("fail")
		require.Error(t, err)
	})
	t.Run("3test", func(t *testing.T) {
		err := r.CheckFileName("")
		require.Error(t, err)
	})

	db.Close()
	t.Run("4test", func(t *testing.T) {
		err := r.CheckFileName("test")
		require.Error(t, err)
	})

}

func TestDB_Create(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	require.NoError(t, err)
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	require.NoError(t, err)
	r.Migrate()

	t.Run("1test", func(t *testing.T) {
		err = r.Create("test")
		require.NoError(t, err)
	})

	t.Run("2test", func(t *testing.T) {
		_, err = db.Exec("DROP TABLE files")
		err = r.Create("")
		if err == errors.New("invalid updated filename") {
			t.Skip()
		}
	})
	t.Run("3test", func(t *testing.T) {
		_, err = db.Exec("DROP TABLE files")
		err = r.Create("test")
	})
}

func TestDB_Migrate(t *testing.T) {
	testfile := "test.db"
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	require.NoError(t, err)
	r := NewSQLiteRepository(db)

	t.Run("1test", func(t *testing.T) {
		err = r.Migrate()
		require.NoError(t, err)
	})
}

func TestDB_Update(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	require.NoError(t, err)
	defer os.Remove(testfile)
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	require.NoError(t, err)
	err = r.Migrate()
	if err != nil {
		return
	}

	t.Run("1test", func(t *testing.T) {
		err = r.Create("fail")
		require.NoError(t, err)
		err = r.Update("test")
		require.NoError(t, err)
	})

	t.Run("2test", func(t *testing.T) {
		err = r.Create("test")
		require.NoError(t, err)
		err = r.Update("test")
		require.NoError(t, err)
	})
	db.Close()
	t.Run("3test", func(t *testing.T) {
		err = r.Update("test")
		require.Error(t, err)
	})

}

func TestSDB_DownloadFileList(t *testing.T) {
	testfile := "test.db"
	fl := file.ListFile{}
	_, err := os.Create(testfile)
	require.NoError(t, err)
	defer os.Remove(testfile)
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	err = r.Migrate()
	require.NoError(t, err)

	t.Run("1test", func(t *testing.T) {
		list, err := r.DownloadFileList()
		require.IsType(t, &fl, list)
		require.NoError(t, err)
	})

	t.Run("2test", func(t *testing.T) {
		err = r.Create("test")
		require.NoError(t, err)
		list, err := r.DownloadFileList()
		require.IsType(t, &fl, list)
		require.NoError(t, err)
	})

	t.Run("3test", func(t *testing.T) {
		err = r.Update("test")
		require.NoError(t, err)
		list, err := r.DownloadFileList()
		require.IsType(t, &fl, list)
		require.NoError(t, err)
	})

	db.Close()
	list, err := r.DownloadFileList()
	require.IsType(t, &fl, list)
	require.Error(t, err)
}

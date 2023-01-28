package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestCurrentTime(t *testing.T) {
	var s string
	require.IsType(t, CurrentTime(), s)
}

func TestNewSQLiteRepository(t *testing.T) {
	var sq *SQLiteRepository
	require.IsType(t, NewSQLiteRepository(DB), sq)

}

func TestSQLiteRepository_All(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}

	_, err = db.Exec("drop table files")
	r.Migrate()
	err = r.Create("test")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Run("1test", func(t *testing.T) {
		r.All()
	})

	t.Run("2test", func(t *testing.T) {
		_, err = db.Exec("drop table files")
		r.All()
	})

}

func TestSQLiteRepository_CheckFileName(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	_, err = db.Exec("drop table files")
	r.Migrate()
	err = r.Create("test")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Run("1test", func(t *testing.T) {
		err := r.CheckFileName("test")
		require.NoError(t, err)
	})
	t.Run("1test", func(t *testing.T) {
		err := r.CheckFileName("fail")
		require.NoError(t, err)
	})

}

func TestSQLiteRepository_Create(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}

	_, err = db.Exec("drop table files")
	r.Migrate()
	err = r.Create("test")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSQLiteRepository_Migrate(t *testing.T) {
	testfile := "test.db"
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	require.NoError(t, err)
	r := NewSQLiteRepository(db)
	err = r.Migrate()
	require.NoError(t, err)

}

func TestSQLiteRepository_Update(t *testing.T) {
	testfile := "test.db"
	_, err := os.Create(testfile)
	defer os.Remove(testfile)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("sqlite3", testfile)
	defer db.Close()
	r := NewSQLiteRepository(db)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}

	_, err = db.Exec("drop table files")
	r.Migrate()
	err = r.Create("test")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = r.Update("test")
	if err != nil {
		t.Fatal(err.Error())
	}

}

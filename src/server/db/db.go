package db

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
)

func InitDb() *sql.DB {
	connInfo := "host=postgres port=5432 user=studier database=studier password=studier sslmode=disable"

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	err = migrateDb(db)
	if err != nil {
		panic(err)
	}

	return db
}

var dbMutex sync.Mutex
var db *sql.DB = InitDb()

type Tx struct {
	Tx *sql.Tx
}

func Begin() (*Tx, error) {
	dbMutex.Lock()

	tx, err := db.Begin()
	if err != nil {
		dbMutex.Unlock()
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

func (tx *Tx) Commit() error {
	defer dbMutex.Unlock()
	err := tx.Tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (tx *Tx) Rollback() error {
	defer dbMutex.Unlock()

	err := tx.Tx.Rollback()
	if err != nil {
		return err
	}

	return nil
}

package dao

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Interface of configureation
type DBConfig interface {
	FormatDSN() string
}

// Prepare sqlx.DB
func initDb(config DBConfig) (*sqlx.DB, error) {
	driverName := "mysql"
	db, err := sqlx.Open(driverName, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open failed: %w", err)
	}

	return db, nil
}

// Transaction handle specific process
// Essentially, it should be abstracted by DB interface
func Transaction(db *sqlx.DB, txFunc func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			log.Println("rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(tx)
	return err
}

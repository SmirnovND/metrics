package db

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/jmoiron/sqlx"
)

func NewDB(c interfaces.ConfigServer) *sqlx.DB {
	dsn := c.GetDBDsn()
	if c.GetDBDsn() == "" {
		dsn = "invalid_dsn"
	}

	db, err := sqlx.Open(
		"postgres",
		dsn,
	)

	if err != nil {
		return db
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)

	fmt.Println("DB connection success!")

	return db
}

package db

import (
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"log"
)

func NewDb(c interfaces.ConfigServer) *sqlx.DB {
	db, err := sqlx.Open(
		"postgres",
		c.GetDbDsn(),
	)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)

	_, err = db.Exec("SET TIME ZONE 'Europe/Moscow';")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connection success!")

	return db
}

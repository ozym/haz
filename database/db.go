package database

import (
	"database/sql"
	"github.com/GeoNet/cfg"
	"log"
	"time"
)

var retry = time.Duration(30) * time.Second

type DB struct {
	*sql.DB
}

func InitPG(c *cfg.DataBase) (DB, error) {
	db, err := sql.Open("postgres", c.Postgres())

	return DB{db}, err
}

func (db *DB) Check() {
	for {
		if err := db.Ping(); err != nil {
			log.Printf("WARN - pinging DB: %s", err)
			log.Println("WARN - waiting then trying DB again.")
			time.Sleep(retry)
			continue
		}
		break
	}
}

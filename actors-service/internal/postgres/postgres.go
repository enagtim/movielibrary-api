package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Db struct {
	*sql.DB
}

func NewConnectDb() (*Db, error) {
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	return &Db{db}, nil
}

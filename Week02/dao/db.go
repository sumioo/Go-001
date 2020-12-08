package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "zclng"
	password = ""
	dbname   = ""
)

func init() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s?sslmode=disable", user, password, host)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	DB = db
}

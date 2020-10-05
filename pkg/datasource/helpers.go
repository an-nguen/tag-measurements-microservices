package datasource

import (
	"database/sql"
	"log"
)

func GetCount(db *sql.DB, tableName string) int {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM $1", tableName)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
		return -1
	}

	return count
}

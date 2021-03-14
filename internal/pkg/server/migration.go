package server

import (
	"database/sql"
	"fmt"
	"github.com/seungha-kim/urlo-go/internal/pkg/util"
	"log"
)

type migration func(*sql.DB)

func RunAllMigrations(db *sql.DB) {
	log.Println("Current database version:", getUserVersion(db))
	for _, m := range migrations {
		m(db)
	}
}

var migrations = []migration{
	migrateFrom0,
}

func migrateFrom0(db *sql.DB) {
	migrate(db, 0, `
	create table urlo (
		id text primary key,
		url text,
		created_at text
	);
	`)
}

func migrate(db *sql.DB, userVersion int, query string) {
	currentVersion := getUserVersion(db)

	if currentVersion != userVersion {
		log.Println("Skip migration for version", userVersion)
		return
	}

	_, err := db.Exec(query)
	util.PanicIfError(err, fmt.Sprintf("Failed to execute query for migration %d", userVersion))

	setUserVersion(db, userVersion+1)
}

func getUserVersion(db *sql.DB) int {
	var result int
	row := db.QueryRow(`pragma user_version;`)
	err := row.Scan(&result)
	util.PanicIfError(err, "Failed to get user_version")
	return result
}

func setUserVersion(db *sql.DB, version int) {
	_, err := db.Exec(fmt.Sprintf("pragma user_version = %d;", version))
	util.PanicIfError(err, "Failed to set user_version")
}

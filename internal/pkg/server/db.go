package server

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/seungha-kim/urlo-go/internal/pkg/conf"
	"github.com/seungha-kim/urlo-go/internal/pkg/util"
	"log"
	"math/rand"
	"time"
)

var db *sql.DB

func init() {
	log.Println("Init db...")
	db = OpenDatabase(conf.DatabasePath)
	RunAllMigrations(db)
}

func OpenDatabase(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	util.PanicIfError(err, "Failed to open database")
	return db
}

func GetUrlById(db *sql.DB, id string) (result string, err error) {
	row := db.QueryRow(`SELECT url FROM urlo WHERE id = ?`, id)
	err = row.Scan(&result)
	return
}

func HasId(db *sql.DB, id string) (result bool, err error) {
	row := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM urlo WHERE id = ?)`, id)
	err = row.Scan(&result)
	return
}

func CreateIdByUrl(db *sql.DB, url string) (string, error) {
	for {
		id := randomId()
		exists, err := HasId(db, id)
		if err != nil {
			return "", err
		}
		if !exists {
			now := time.Now().Format(time.RFC3339)
			_, err := db.Exec(`INSERT INTO urlo VALUES (?, ?, ?)`, id, url, now)
			return id, err
		}
	}
}

var idCandidates = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomId() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var runes [6]rune
	for i := 0; i < len(runes); i++ {
		runes[i] = idCandidates[r.Intn(len(idCandidates))]
	}
	return string(runes[:])
}

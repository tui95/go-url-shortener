package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tui95/go-url-shortener/internal/config"
)

func NewDB(url string) *sql.DB {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping db: %v\n", err)
	}
	return db
}

func CreateTableIfNotExists(db *sql.DB) {
	dat, err := os.ReadFile(filepath.Join(config.BASE_DIR, "internal/database/migrations/001_create_url_mapping.sql"))
	if err != nil {
		log.Fatalf("Unable to open create url mapping sql file: %v\n", err)
	}
	_, err = db.Exec(string(dat))
	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}
}

func CreateUrlMapping(db *sql.DB, url string) (int, error) {
	statement, err := db.Prepare("INSERT INTO url_mapping(url) VALUES(?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v\n", err)
		return -1, err
	}
	result, err := statement.Exec(url)
	if err != nil {
		log.Printf("Failed to exec statement: %v\n", err)
		return -1, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func GetUrlById(db *sql.DB, id int) (string, error) {
	var url string
	row := db.QueryRow("SELECT url from url_mapping where id = ?", id)
	err := row.Scan(&url)
	if err != nil {
		log.Printf("Unable to get url by id: %v", err)
		return "", err
	}
	return url, nil
}

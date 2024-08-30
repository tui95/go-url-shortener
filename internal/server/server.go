package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tui95/go-url-shortener/internal/config"
	"github.com/tui95/go-url-shortener/internal/database"
)

func RunServer(port string) {
	db := database.NewDB(filepath.Join(config.BASE_DIR, "db.sqlite3"))
	defer db.Close()
	database.CreateTableIfNotExists(db)
	router := NewRouter(db)
	fmt.Printf("Running on %v\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

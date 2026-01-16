package migration

import (
	"database/sql"
	"log"
	"time"

	"github.com/pressly/goose/v3"
)

const (
	attempts = 5
)

func RunMigrations(db *sql.DB, migrationsDir string) {

	time.Sleep(5 * time.Second)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("goose set dialect: %v", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Printf("goose up failed: %v", err)
	}

	log.Println("migrations applied successfully")
}

func RollbackLast(db *sql.DB, migrationsDir string) {
	if err := goose.Down(db, migrationsDir); err != nil {
		log.Fatalf("goose down failed: %v", err)
	}
	log.Println("last migration rolled back successfully!")
}

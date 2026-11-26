package main

import (
	"log"

	"karsa/internal/config"
	"karsa/internal/database"
	"karsa/internal/seed"
)

func main() {
	cfg := config.Load()
	if err := database.Migrate(cfg.DBDSN, "migrations"); err != nil {
		log.Fatal("migrasi gagal: ", err)
	}
	db, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatal("koneksi database gagal: ", err)
	}
	if err := seed.Run(db, cfg); err != nil {
		log.Fatal("seed gagal: ", err)
	}
}

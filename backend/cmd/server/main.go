package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata"

	"karsa/internal/config"
	"karsa/internal/database"
	"karsa/internal/mailer"
	"karsa/internal/routes"
	"karsa/internal/scheduler"
	"karsa/internal/seed"
	"karsa/internal/worker"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg := config.Load()
	if cfg.DBDSN == "" || cfg.JWTSecret == "" {
		log.Fatal("DB_DSN dan JWT_SECRET wajib diisi di .env")
	}
	if err := database.Migrate(cfg.DBDSN, "migrations"); err != nil {
		log.Fatal("migrasi gagal: ", err)
	}
	db, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatal("koneksi database gagal: ", err)
	}
	if cfg.RunSeed {
		if err := seed.Run(db, cfg); err != nil {
			log.Fatal("seed gagal: ", err)
		}
	}
	router, err := routes.New(db, cfg)
	if err != nil {
		log.Fatal("router gagal: ", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var background sync.WaitGroup
	background.Go(func() { worker.Run(ctx, db, mailer.Mailer{Config: cfg}) })
	background.Go(func() { scheduler.Run(ctx, db) })

	server := &http.Server{Addr: ":" + cfg.Port, Handler: router, ReadHeaderTimeout: shutdownTimeout}
	go func() {
		log.Println("server berjalan di port", cfg.Port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server berhenti: ", err)
		}
	}()

	<-ctx.Done()
	log.Println("mematikan server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown:", err)
	}
	background.Wait()
}

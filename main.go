package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JscorpTech/paymento/database"
	"github.com/JscorpTech/paymento/handlers"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	WatchBotID int64 = 7612977626
)

func main() {
	// Load env
	if err := godotenv.Load(".env"); err != nil {
		panic(".env file not loaded: " + err.Error())
	}

	// DB
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	database.InitTables(db)
	// Logger
	log, _ := zap.NewDevelopment(
		zap.IncreaseLevel(zapcore.InfoLevel),
		zap.AddStacktrace(zapcore.FatalLevel),
	)
	defer log.Sync()

	handler := handlers.NewHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("/create/transaction/", handler.HandlerHome)

	srv := &http.Server{
		Addr:    ":8084",
		Handler: mux,
	}

	// Context with signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	log.Info("server started", zap.String("addr", srv.Addr))

	<-ctx.Done()
	log.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", zap.Error(err))
	} else {
		log.Info("server exited properly")
	}
}

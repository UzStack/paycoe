package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JscorpTech/paymento/internal/config"
	"github.com/JscorpTech/paymento/internal/domain"
	"github.com/JscorpTech/paymento/internal/http/routes"
	"github.com/JscorpTech/paymento/internal/infra"
	"github.com/JscorpTech/paymento/internal/repository"
	"github.com/JscorpTech/paymento/internal/usecase"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func printHelp() {
	fmt.Println("Usage paycue [options]")
	fmt.Println("Options:")
	fmt.Println("  -h                 Show this help message")
	fmt.Println("  --telegram         Telegram accountni ulash")
	os.Exit(0)
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(".env file not loaded: " + err.Error())
	}
	log, _ := zap.NewDevelopment(
		zap.IncreaseLevel(zapcore.InfoLevel),
		zap.AddStacktrace(zapcore.FatalLevel),
	)
	defer log.Sync()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h":
			printHelp()
		case "--telegram":
			if err := infra.Mtproto(context.Background(), nil, zap.NewNop(), 0, false, nil, nil); err != nil {
				panic(err)
			}
			fmt.Println("")
			fmt.Println("Account qo'shildi")
			os.Exit(0)
		}
	}

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.InitTables(db)
	tasks := make(chan domain.Task, 10)

	mux := http.NewServeMux()
	routes.InitRoutes(mux, db, log, tasks, cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := usecase.InitWorker(ctx, log, tasks, cfg); err != nil {
		log.Error("worker init failed", zap.Any("error", err.Error()))
	}
	defer close(tasks)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()
	go func() {
		if err := infra.Mtproto(ctx, db, log, cfg.WatchID, true, tasks, cfg); err != nil {
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

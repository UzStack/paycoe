package routes

import (
	"database/sql"
	"net/http"

	"github.com/JscorpTech/paymento/internal/config"
	"github.com/JscorpTech/paymento/internal/domain"
	"github.com/JscorpTech/paymento/internal/http/handlers"
	"go.uber.org/zap"
)

func InitRoutes(mux *http.ServeMux, db *sql.DB, log *zap.Logger, tasks chan domain.Task, cfg *config.Config) {
	handler := handlers.NewHandler(db, log, tasks, cfg)
	mux.HandleFunc("/create/transaction/", handler.HandlerHome)
}

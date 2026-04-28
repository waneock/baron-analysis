package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
	"skinbaron-analyzer/pkg/logger"
	"skinbaron-analyzer/services/reporting/internal/client/parsinggrpc"
	"skinbaron-analyzer/services/reporting/internal/config"
	httphndl "skinbaron-analyzer/services/reporting/internal/transport/http"
	"skinbaron-analyzer/services/reporting/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	// logger
	log := logger.MustLoad(cfg.Env)

	// database
	dbCfg := makeDBConfigData(cfg)

	db, err := db.New(*dbCfg)
	if err != nil {
		log.Error("error when trying to create a new database",
			"error", err)
		os.Exit(1)
	}

	defer db.Close()

	parsingClient, err := parsinggrpc.New("parsing:50051")
	if err != nil {
		log.Error("cannot create parsing client",
			"error", err)
		os.Exit(1)
	}

	listOffersUC := usecase.NewListOffers(parsingClient)
	syncOffersUC := usecase.NewSyncOffers(parsingClient)

	httpHandler := httphndl.NewOffersHandler(syncOffersUC, listOffersUC)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/offers", func(r chi.Router) {
			r.Get("/list", httpHandler.ListOffers)
			r.Post("/sync", httpHandler.SyncOffers)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	fmt.Println("run server on: ", server.Addr)

	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server stopped",
			"error", err)
		os.Exit(1)
	}
}

func makeDBConfigData(cfg *config.Config) *db.DBConfigData {
	return &db.DBConfigData{
		ConnUrl:         env.GetDBUrl(),
		ConnTimeout:     cfg.DBConfig.ConnTimeout,
		ConnMaxIdleTime: cfg.DBConfig.ConnMaxIdleTime,
		ConnMaxLifeTime: cfg.DBConfig.ConnMaxLifeTime,
		MaxIdleConns:    cfg.DBConfig.MaxIdleConns,
		MaxOpenConns:    cfg.DBConfig.MaxOpenConns,
	}
}

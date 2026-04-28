package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
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

	parsingClient, err := parsinggrpc.New("parsing:50051")
	if err != nil {
		log.Fatal("cannot create parsing client: ", err)
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

	if err = server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server stopped", "err", err)
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

package main

import (
	"log"
	"net/http"

	"github.com/nextpricedevelopers/go-next/internal/config"
	"github.com/nextpricedevelopers/go-next/internal/config/logger"
	"github.com/nextpricedevelopers/go-next/internal/handler/payment"
	"github.com/nextpricedevelopers/go-next/pkg/adapter/mongodb"
	"github.com/nextpricedevelopers/go-next/pkg/server"
	payment_service "github.com/nextpricedevelopers/go-next/pkg/service/payment"

	"github.com/go-chi/chi/v5"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {
	logger.Info("About to start user application")
	conf := config.NewConfig()
	mdb_conn := mongodb.New(conf)

	r := chi.NewRouter()

	r.Get("/", healthcheck)
	payment_service := payment_service.NewPaymentService(mdb_conn)
	payment.RegisterPaymentAPIHandlers(r, payment_service)

	srv := server.NewHTTPServer(r, conf)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server Run on [Port: %s], [Mode: %s], [Version: %s], [Commit: %s]", conf.PORT, conf.Mode, VERSION, COMMIT)

	select {}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"MSG": "Server Ok", "codigo": 200}`))
}

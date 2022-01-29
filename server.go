package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ShauryaAg/ProductAPI/models/db"
	"github.com/ShauryaAg/ProductAPI/routes"
	"github.com/rs/cors"
)

func main() {
	r := routes.GetRoutes()

	ctx := context.Background()
	client, err := db.InitDatabase("mongo", ctx)
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal("errrrr", err)
	}

	db.DBCon = client.Database("mongo")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Accept", "Accept-Language", "Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      c.Handler(r),
	}

	log.Fatal(srv.ListenAndServe())
}

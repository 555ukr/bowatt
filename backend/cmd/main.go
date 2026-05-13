package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/555ukr/bowatt/internal/api"
	"github.com/555ukr/bowatt/pkg/database"
	"github.com/555ukr/bowatt/pkg/storage"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	store := storage.NewLocalStorageService("./uploads")
	os.MkdirAll("./uploads", 0o755)

	// Connect to PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://insta:insta_secret@localhost:5432/insta_like?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := database.NewPostgresPhotoRepository(db)

	router := mux.NewRouter()
	router.Use(api.LoggingMiddleware)
	router.HandleFunc("/health", api.HealthHandler).Methods("GET")
	router.HandleFunc("/upload", api.UploadHandler(store, repo)).Methods("POST")
	http.Handle("/", router)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Println("web server is about start")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

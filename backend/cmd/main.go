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
	"github.com/555ukr/bowatt/internal/websocket"
	"github.com/555ukr/bowatt/pkg/database"
	"github.com/555ukr/bowatt/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("[ERROR]: no .env file found, using environment variables")
	}

	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		log.Fatal("[ERROR]: UPLOAD_PATH is not set")
	}
	store := storage.NewLocalStorageService(uploadPath)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("[ERROR]: failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := database.NewPostgresPhotoRepository(db)

	hub := websocket.NewHub()

	router := mux.NewRouter()
	router.Use(api.CORSMiddleware)
	router.Use(api.LoggingMiddleware)
	router.HandleFunc("/health", api.HealthHandler).Methods("GET")
	router.HandleFunc("/upload", api.UploadHandler(store, repo, hub)).Methods("POST")
	router.HandleFunc("/photos", api.ListPhotosHandler(store, repo)).Methods("GET")
	router.HandleFunc("/ws", hub.Handler())

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8000"
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Println("[INFO]: web server is about start")
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
	log.Println("[INFO]: shutting down")
	os.Exit(0)
}

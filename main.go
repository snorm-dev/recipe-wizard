package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Could not load .env file: ", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Println("Could not load custom port. Using default port 8080")
		port = "8080"
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	v1 := chi.NewRouter()
	r.Mount("/v1", v1)

	v1.Get("/ping", handlePing)

	server := &http.Server{
		Addr:              "localhost:" + port,
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("Listening on port: ", port)
	log.Fatal(server.ListenAndServe())
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong!"))
}

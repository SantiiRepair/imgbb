package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/santiirepair/imgbb/uploader"
)

func main() {
	defaultAPIKey := os.Getenv("IMGBB_API_KEY")
	if defaultAPIKey == "" {
		log.Println("WARNING: IMGBB_API_KEY is not configured in the environment variables. It will be required in every request.")
	}

	maxSizeStr := os.Getenv("MAX_UPLOAD_SIZE_MB")
	maxSizeMB := 10 // default
	if s, err := strconv.Atoi(maxSizeStr); err == nil && s > 0 {
		maxSizeMB = s
	}
	maxBytes := int64(maxSizeMB) << 20

	mux := http.NewServeMux()

	// Health Check Endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Limit upload size
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		if err := r.ParseMultipartForm(maxBytes); err != nil {
			http.Error(w, "File too large or malformed", http.StatusBadRequest)
			return
		}

		// Determine which API key to use
		apiKey := r.FormValue("api_key")
		if apiKey == "" {
			apiKey = defaultAPIKey
		}

		if apiKey == "" {
			http.Error(w, "API key not provided in the request or configured in the server", http.StatusUnauthorized)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "The 'image' field is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading the image", http.StatusInternalServerError)
			return
		}

		// Context for the upload request
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		// Upload to ImgBB directly
		url, err := uploader.Upload(ctx, apiKey, data)
		if err != nil {
			log.Printf("Error uploading to ImgBB: %v", err)
			http.Error(w, "Error uploading the image", http.StatusInternalServerError)
			return
		}

		// Respond with the URL
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"url": url,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("ImgBB upload server started on port :%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

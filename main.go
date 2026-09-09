package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/santiirepair/imgbb/uploader"
)

func main() {
	defaultAPIKey := os.Getenv("IMGBB_API_KEY")
	if defaultAPIKey == "" {
		log.Println("WARNING: IMGBB_API_KEY is not configured in the environment variables. It will be required in every request.")
	}

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Limit upload size to 10MB
		r.ParseMultipartForm(10 << 20)

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

		// Upload to ImgBB directly
		url, err := uploader.Upload(apiKey, data)
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
	log.Printf("ImgBB upload server started on port :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

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

	app := fiber.New(fiber.Config{
		BodyLimit: maxSizeMB * 1024 * 1024,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST, OPTIONS, GET",
	}))

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Upload Endpoint
	app.Post("/upload", func(c *fiber.Ctx) error {
		// Determine which API key to use
		apiKey := c.FormValue("api_key")
		if apiKey == "" {
			apiKey = defaultAPIKey
		}

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key not provided in the request or configured in the server",
			})
		}

		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "The 'image' field is required",
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to open image",
			})
		}
		defer file.Close()

		data := make([]byte, fileHeader.Size)
		if _, err := file.Read(data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read image",
			})
		}

		// Context for the upload request
		ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
		defer cancel()

		// Upload to ImgBB
		url, err := uploader.Upload(ctx, apiKey, data)
		if err != nil {
			log.Printf("Error uploading to ImgBB: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error uploading the image to ImgBB",
			})
		}

		return c.JSON(fiber.Map{"url": url})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown channel
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server gracefully...")

		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	log.Printf("ImgBB upload server started on port :%s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}

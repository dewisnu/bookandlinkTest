package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Initialize necessary directories
	err := initializeDirectories()
	if err != nil {
		log.Fatalf("Failed to initialize directories: %v", err)
	}

	// Initialize database connection
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize RabbitMQ connection
	rabbitmq, err := initRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define API routes
	router.POST("/upload", handleUpload(db, rabbitmq))
	router.GET("/jobs", getJobs(db))
	router.GET("/jobs/:id", getJob(db))
	router.GET("/jobs/status/:status", getJobsByStatus(db))
	router.POST("/jobs/:id/retry", retryJob(db, rabbitmq))
	router.GET("/images-uploaded/:filename", serveImageUploaded())
	router.GET("/images-compressed/:filename", serveImageCompressed())
	router.POST("/compressed", handleCompressedUpload())
	router.Static("/compressed", "./compressed")

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Publisher service starting on port %s...", port)
	router.Run(fmt.Sprintf(":%s", port))
}

// initializeDirectories creates the necessary directories for the service
func initializeDirectories() error {
	dirs := []string{
		"./uploads",
		"./compressed",
	}

	for _, dir := range dirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", dir, err)
		}

		err = os.MkdirAll(absPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", absPath, err)
		}
		log.Printf("Ensured directory exists: %s", absPath)
	}

	return nil
}

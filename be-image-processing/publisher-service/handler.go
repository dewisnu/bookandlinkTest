package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// handleUpload handles image upload requests
func handleUpload(db *sql.DB, rabbitmq *RabbitMQ) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get uploaded files
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		files := form.File["images"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No images uploaded"})
			return
		}

		// Process each file
		var jobIDs []int64
		for _, file := range files {
			// Check if file is an image
			if !isImage(file.Filename) {
				log.Printf("Skipping non-image file: %s", file.Filename)
				continue
			}

			// Generate unique filename
			filename := generateUniqueFilename(file.Filename)
			filepath := filepath.Join("uploads", filename)

			// Save file to disk
			if err := c.SaveUploadedFile(file, filepath); err != nil {
				log.Printf("Error saving file %s: %v", filename, err)
				continue
			}

			// Get file size
			fileInfo, err := os.Stat(filepath)
			if err != nil {
				log.Printf("Error getting file info for %s: %v", filename, err)
				os.Remove(filepath) // Clean up file
				continue
			}

			originalSize := fileInfo.Size()

			// Create job in database
			jobID, err := createImageJob(db, filename, originalSize)
			if err != nil {
				log.Printf("Error creating job for %s: %v", filename, err)
				os.Remove(filepath) // Clean up file
				continue
			}

			// Publish job to queue
			err = rabbitmq.PublishJob(jobID, filename)
			if err != nil {
				log.Printf("Error publishing job for %s: %v", filename, err)
				updateJobStatus(db, jobID, "failed")
				continue
			}

			jobIDs = append(jobIDs, jobID)
			log.Printf("Successfully processed upload: %s, size: %d bytes, job ID: %d", filename, originalSize, jobID)
		}

		if len(jobIDs) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process any uploaded images"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully processed %d images", len(jobIDs)),
			"job_ids": jobIDs,
		})
	}
}

// handleCompressedUpload handles receiving compressed images from worker service
func handleCompressedUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse form data
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing compressed image"})
			return
		}

		// Optional: validate if it's an image (e.g., jpg/png)
		if !isImage(file.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not a valid image"})
			return
		}

		// Save to ./compressed directory
		filename := filepath.Base(file.Filename)
		savePath := filepath.Join("compressed", filename)

		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			log.Printf("Failed to save compressed file %s: %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save compressed file"})
			return
		}

		log.Printf("Received and saved compressed image: %s", filename)
		c.JSON(http.StatusOK, gin.H{
			"message":  "Compressed image received",
			"filename": filename,
			"path":     "/compressed/" + filename,
		})
	}
}

// getJobs returns all image processing jobs
func getJobs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobs, err := getImageJobs(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching jobs: %v", err)})
			return
		}

		c.JSON(http.StatusOK, jobs)
	}
}

// getJob returns a single job by ID
func getJob(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
			return
		}

		job, err := getImageJob(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching job: %v", err)})
			}
			return
		}

		c.JSON(http.StatusOK, job)
	}
}

// getJobsByStatus returns jobs filtered by status
func getJobsByStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Param("status")
		// Validate status
		if status != "pending" && status != "processing" && status != "completed" && status != "failed" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: pending, processing, complete, failed"})
			return
		}

		jobs, err := getImageJobsByStatus(db, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching jobs: %v", err)})
			return
		}

		c.JSON(http.StatusOK, jobs)
	}
}

// retryJob retries a failed job
func retryJob(db *sql.DB, rabbitmq *RabbitMQ) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
			return
		}

		// Get job from database
		job, err := getImageJob(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching job: %v", err)})
			}
			return
		}

		// Only failed jobs can be retried
		if job.Status != "failed" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only failed jobs can be retried"})
			return
		}

		err = updateJobStatus(db, id, "pending")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating job status: %v", err)})
			return
		}

		// Publish job to queue
		err = rabbitmq.PublishJob(id, job.Filename)
		if err != nil {
			// If publishing fails, revert status to failed
			updateJobStatus(db, id, "failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error publishing job to queue: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Job queued for retry"})
	}
}

// serveImageUploaded serves uploaded images from disk
func serveImageUploaded() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")

		// Validate filename to avoid path traversal
		if containsPathTraversal(filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
			return
		}

		// Correct path to the file in container
		imagePath := filepath.Join("/app/uploads", filename)

		// Check if the file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}

		// Serve the file
		c.File(imagePath)
	}
}

func serveImageCompressed() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")

		// For security, ensure filename doesn't contain path traversal
		if containsPathTraversal(filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
			return
		}

		// Path to compressed image
		imagePath := filepath.Join("/app/compressed", filename)

		// Check if file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}

		// Serve the file
		c.File(imagePath)
	}
}

// isImage checks if a filename has an image extension
func isImage(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}

// Slugify makes the filename lowercase, alphanumeric, and hyphen-separated
func slugify(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)
	// Replace spaces with dashes
	name = strings.ReplaceAll(name, " ", "-")
	// Remove any character that's not alphanumeric or dash
	reg := regexp.MustCompile("[^a-z0-9\\-]+")
	name = reg.ReplaceAllString(name, "")
	return name
}

// generateUniqueFilename generates a clean, unique filename
func generateUniqueFilename(original string) string {
	ext := filepath.Ext(original)
	name := strings.TrimSuffix(original, ext)
	slug := slugify(name)
	timestamp := time.Now().Unix() // or use uuid.New().String() for full uniqueness
	return fmt.Sprintf("%s-%d%s", slug, timestamp, ext)
}

// getTimestamp returns current Unix timestamp
func getTimestamp() int64 {
	return int64(getRandTimestamp())
}

func getRandTimestamp() int64 {
	return 123456789 // This is a placeholder for the actual implementation
}

// containsPathTraversal checks if a filename contains path traversal patterns
func containsPathTraversal(filename string) bool {
	if filename == "." || filename == ".." ||
		filepath.IsAbs(filename) ||
		filepath.Clean(filename) != filename {
		return true
	}
	return false
}

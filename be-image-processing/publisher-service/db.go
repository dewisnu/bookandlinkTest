package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type ImageJob struct {
	ID                 int64      `json:"id"`
	Filename           string     `json:"filename"`
	OriginalSize       *int64     `json:"original_size"`
	CompressedSize     *int64     `json:"compressed_size"`
	CompressedFileName *string    `json:"compressed_file_name"`
	Status             string     `json:"status"`
	ErrorMessage       *string    `json:"error_message"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}

func InitDB() (*sql.DB, error) {
	// Get database connection parameters from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "image_processor")

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to database")
	return db, nil
}

// createImageJob inserts a new job into the database
//func createImageJob(db *sql.DB, filename string, originalSize int64) (int64, error) {
//	query := `
//		INSERT INTO image_jobs (filename, original_size, status)
//		VALUES ($1, $2, 'pending')
//		RETURNING id
//	`
//
//	var id int64
//	err := db.QueryRow(query, filename, originalSize).Scan(&id)
//	if err != nil {
//		return 0, fmt.Errorf("error creating image job: %w", err)
//	}
//
//	return id, nil
//}

// getImageJobs retrieves all image jobs from the database
func getImageJobs(db *sql.DB) ([]ImageJob, error) {
	query := `
		SELECT id, filename, original_size, compressed_size, compressed_file_name, 
			   status, error_message, created_at, updated_at
		FROM image_jobs
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying image jobs: %w", err)
	}
	defer rows.Close()

	var jobs []ImageJob
	for rows.Next() {
		var job ImageJob
		err := rows.Scan(
			&job.ID, &job.Filename, &job.OriginalSize, &job.CompressedSize,
			&job.CompressedFileName, &job.Status, &job.ErrorMessage,
			&job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning image job row: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating image job rows: %w", err)
	}

	return jobs, nil
}

// getImageJobsByStatus retrieves image jobs filtered by status
//func getImageJobsByStatus(db *sql.DB, status string) ([]ImageJob, error) {
//	query := `
//		SELECT id, filename, original_size, compressed_size, compressed_file_name,
//			   status, error_message, created_at, updated_at
//		FROM image_jobs
//		WHERE status = $1
//		ORDER BY created_at DESC
//		LIMIT 100
//	`
//
//	rows, err := db.Query(query, status)
//	if err != nil {
//		return nil, fmt.Errorf("error querying image jobs by status: %w", err)
//	}
//	defer rows.Close()
//
//	var jobs []ImageJob
//	for rows.Next() {
//		var job ImageJob
//		err := rows.Scan(
//			&job.ID, &job.Filename, &job.OriginalSize, &job.CompressedSize,
//			&job.CompressedFileName, &job.Status, &job.ErrorMessage,
//			&job.CreatedAt, &job.UpdatedAt,
//		)
//		if err != nil {
//			return nil, fmt.Errorf("error scanning image job row: %w", err)
//		}
//		jobs = append(jobs, job)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, fmt.Errorf("error iterating image job rows: %w", err)
//	}
//
//	return jobs, nil
//}

// getImageJob retrieves a single image job by ID
//func getImageJob(db *sql.DB, id int64) (ImageJob, error) {
//	query := `
//		SELECT id, filename, original_size, compressed_size, compressed_file_name,
//			   status, error_message, created_at, updated_at
//		FROM image_jobs
//		WHERE id = $1
//	`
//
//	var job ImageJob
//	err := db.QueryRow(query, id).Scan(
//		&job.ID, &job.Filename, &job.OriginalSize, &job.CompressedSize,
//		&job.CompressedFileName, &job.Status, &job.ErrorMessage,
//		&job.CreatedAt, &job.UpdatedAt,
//	)
//	if err != nil {
//		return ImageJob{}, fmt.Errorf("error getting image job: %w", err)
//	}
//
//	return job, nil
//}

// updateJobStatus updates the status of a job
//func updateJobStatus(db *sql.DB, id int64, status string) error {
//	query := `
//		UPDATE image_jobs
//		SET status = $2
//		WHERE id = $1
//	`
//
//	_, err := db.Exec(query, id, status)
//	if err != nil {
//		return fmt.Errorf("error updating job status: %w", err)
//	}
//
//	return nil
//}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

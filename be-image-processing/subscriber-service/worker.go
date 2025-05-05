package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func processJob(db *sql.DB, msg amqp.Delivery) {
	var jobMsg JobMessage
	if err := json.Unmarshal(msg.Body, &jobMsg); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		msg.Nack(false, false)
		return
	}

	// Cek retry count
	retryCount := getRetryCount(msg.Headers)
	if retryCount >= 3 {
		log.Printf("Job ID %d reached max retries. Sending to DLQ.", jobMsg.ID)
		updateJobStatus(db, jobMsg.ID, "failed", "max retry reached", nil, nil)
		msg.Reject(false) // go to DLQ
		return
	}

	updateJobStatus(db, jobMsg.ID, "processing", "", nil, nil)

	query := `SELECT filename FROM image_jobs WHERE id = $1`
	var filename string
	err := db.QueryRow(query, jobMsg.ID).Scan(&filename)
	if err != nil {
		log.Printf("Error fetching job data: %v", err)
		updateJobStatus(db, jobMsg.ID, "failed", fmt.Sprintf("DB error: %v", err), nil, nil)
		msg.Nack(false, false)
		return
	}

	// Download image from provider service
	imageURL := "http://publisher-service:8080/images-uploaded/" + filename
	resp, err := http.Get(imageURL)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Error downloading image: %v", err)
		updateJobStatus(db, jobMsg.ID, "failed", "failed to fetch image from provider", nil, nil)
		msg.Nack(false, false)
		return
	}
	defer resp.Body.Close()

	tempInput := filepath.Join(os.TempDir(), filename)
	out, err := os.Create(tempInput)
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		msg.Nack(false, false)
		return
	}
	io.Copy(out, resp.Body)
	out.Close()

	outputPath := filepath.Join("./compressed", "compressed_"+filename)

	compressedSize, err := compressImage(tempInput, outputPath)
	if err != nil {
		log.Printf("Error compressing image: %v", err)
		updateJobStatus(db, jobMsg.ID, "failed", fmt.Sprintf("Compression error: %v", err), nil, nil)
		msg.Nack(false, false)
		return
	}

	fileData, err := os.Open(outputPath)
	if err != nil {
		log.Printf("Failed to open compressed file: %v", err)
		msg.Nack(false, false)
		return
	}
	defer fileData.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(outputPath))
	if err != nil {
		log.Printf("Error creating multipart: %v", err)
		msg.Nack(false, false)
		return
	}
	io.Copy(part, fileData)
	writer.Close()

	uploadResp, err := http.Post("http://publisher-service:8080/compressed", writer.FormDataContentType(), body)
	if err != nil || uploadResp.StatusCode != 200 {
		log.Printf("Error uploading compressed image: %v", err)
		msg.Nack(false, false)
		return
	}

	compressedFileName := "compressed_" + filename
	updateJobStatus(db, jobMsg.ID, "completed", "", &compressedSize, &compressedFileName)

	log.Printf("Successfully processed job %d", jobMsg.ID)
	msg.Ack(false)
}

func updateJobStatus(db *sql.DB, id int, status, errorMsg string, compressedSize *int64, compressedURL *string) {
	query := `
        UPDATE image_jobs 
        SET status = $1, error_message = $2, compressed_size = $3, compressed_file_name = $4, updated_at = NOW() 
        WHERE id = $5
    `
	_, err := db.Exec(query, status, errorMsg, compressedSize, compressedURL, id)
	if err != nil {
		log.Printf("Failed to update job status: %v", err)
	}
}

func getRetryCount(headers amqp.Table) int {
	xDeathRaw, ok := headers["x-death"]
	if !ok {
		return 0
	}

	xDeathList, ok := xDeathRaw.([]interface{})
	if !ok || len(xDeathList) == 0 {
		return 0
	}

	xDeath, ok := xDeathList[0].(amqp.Table)
	if !ok {
		return 0
	}

	count, ok := xDeath["count"].(int64)
	if !ok {
		return 0
	}

	return int(count)
}

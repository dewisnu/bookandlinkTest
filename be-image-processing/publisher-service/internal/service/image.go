package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"publisher-service/internal/util/helper"
	"publisher-service/pkg/dto"
)

func (s *service) HandleUpload(g *gin.Context, files []*multipart.FileHeader) (imageResponse dto.ImageResponse, err error) {
	var jobIDs []int64
	for _, file := range files {
		if !helper.IsImage(file.Filename) {
			slog.Warn(fmt.Sprintf("Skipping non-image file: %s", file.Filename))
			continue
		}

		// Generate unique filename
		filename := helper.GenerateUniqueFilename(file.Filename)
		filepathImage := filepath.Join("uploads", filename)

		// Save file to disk
		if err := g.SaveUploadedFile(file, filepathImage); err != nil {
			log.Printf("Error saving file %s: %v", filename, err)
			continue
		}

		// Get file size
		fileInfo, err := os.Stat(filepathImage)
		if err != nil {
			slog.Error(fmt.Sprintf("Error getting file info for %s: %v", filename, err))
			os.Remove(filepathImage) // Clean up file
			continue
		}

		originalSize := fileInfo.Size()

		// Create job in database
		jobID, err := s.repository.CreateImageJob(filename, originalSize)
		if err != nil {
			slog.Error(fmt.Sprintf("Error creating job for %s: %v", filename, err))
			os.Remove(filepathImage) // Clean up file
			continue
		}

		// Publish job to queue
		err = s.rabbitmq.PublishJob(jobID, filename)
		if err != nil {
			slog.Error(fmt.Sprintf("Error publishing job for %s: %v", filename, err))
			s.repository.UpdateJobStatus(jobID, "failed")
			continue
		}

		err = s.repository.UpdateJobStatus(jobID, "processing")
		if err != nil {
			slog.Error(fmt.Sprintf("Error updating job status: %v", err))
			continue
		}

		jobIDs = append(jobIDs, jobID)
		slog.Info(fmt.Sprintf("Successfully processed upload: %s, size: %d bytes, job ID: %d", filename, originalSize, jobID))
	}
	return dto.ImageResponse{JobIDs: jobIDs}, nil
}

func (s *service) ServeImageUploaded(filename string) (imagePath string, isExist bool, err error) {
	imagePath = filepath.Join("/app/uploads", filename)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		isExist = os.IsNotExist(err)
		slog.Error("Image not found")
		return imagePath, isExist, err
	}
	return
}

func (s *service) ServeImageCompressed(filename string) (imagePath string, isExist bool, err error) {
	imagePath = filepath.Join("/app/compressed", filename)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		isExist = os.IsNotExist(err)
		slog.Error("Image not found")
		return imagePath, isExist, err
	}
	return
}

func (s *service) CompressedUpload(g *gin.Context, file *multipart.FileHeader) (compressedImageResponse dto.CompressedImageResponse, err error) {
	filename := filepath.Base(file.Filename)
	savePath := filepath.Join("compressed", filename)

	err = g.SaveUploadedFile(file, savePath)
	if err != nil {
		log.Printf("Failed to save compressed file %s: %v", filename, err)
		slog.Error(fmt.Sprintf("Failed to save compressed file %s: %v", filename, err))
		return dto.CompressedImageResponse{}, fmt.Errorf("failed to save compressed file")
	}

	slog.Info(fmt.Sprintf("Received and saved compressed image: %s", filename))

	return dto.CompressedImageResponse{
		Filename: filename,
		Path:     "/compressed/" + filename,
	}, nil
}

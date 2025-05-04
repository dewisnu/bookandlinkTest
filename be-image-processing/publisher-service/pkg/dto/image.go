package dto

import "time"

type ImageResponse struct {
	JobIDs []int64 `json:"imageId"`
}

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

type CompressedImageResponse struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

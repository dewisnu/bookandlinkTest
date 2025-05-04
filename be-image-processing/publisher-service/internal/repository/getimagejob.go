package repository

import (
	"fmt"
	"publisher-service/pkg/dto"
)

func (r repository) GetImageJob(id int64) (dto.ImageJob, error) {
	query := `
		SELECT id, filename, original_size, compressed_size, compressed_file_name,
			   status, error_message, created_at, updated_at
		FROM image_jobs
		WHERE id = $1
	`

	var job dto.ImageJob
	err := r.db.QueryRow(query, id).Scan(
		&job.ID, &job.Filename, &job.OriginalSize, &job.CompressedSize,
		&job.CompressedFileName, &job.Status, &job.ErrorMessage,
		&job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		return dto.ImageJob{}, fmt.Errorf("error getting image job: %w", err)
	}

	return job, nil
}

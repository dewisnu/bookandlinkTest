package repository

import (
	"fmt"
	"publisher-service/pkg/dto"
)

func (r repository) GetImageJobsByStatus(status string) ([]dto.ImageJob, error) {
	query := `
		SELECT id, filename, original_size, compressed_size, compressed_file_name,
			   status, error_message, created_at, updated_at
		FROM image_jobs
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("error querying image jobs by status: %w", err)
	}
	defer rows.Close()

	var jobs []dto.ImageJob
	for rows.Next() {
		var job dto.ImageJob
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

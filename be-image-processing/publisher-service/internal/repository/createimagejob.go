package repository

import (
	"fmt"
)

func (r repository) CreateImageJob(filename string, originalSize int64) (int64, error) {
	query := `
		INSERT INTO image_jobs (filename, original_size, status)
		VALUES ($1, $2, 'pending')
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(query, filename, originalSize).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating image job: %w", err)
	}

	return id, nil
}

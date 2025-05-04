package repository

import (
	"fmt"
)

func (r repository) UpdateJobStatus(id int64, status string) error {
	query := `
		UPDATE image_jobs
		SET status = $2
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id, status)
	if err != nil {
		return fmt.Errorf("error updating job status: %w", err)
	}

	return nil
}

package service

import (
	"errors"
	"fmt"
	"log/slog"
	"publisher-service/pkg/dto"
)

func (s *service) GetJobs() (imageJobsResponse []dto.ImageJob, err error) {
	jobs, err := s.repository.GetImageJobs()
	if err != nil {
		slog.Error(fmt.Sprintf("Error fetching jobs: %v", err))
		return nil, err
	}

	return jobs, nil
}

func (s *service) GetJob(id int64) (imageJobResponse dto.ImageJob, err error) {
	imageJobResponse, err = s.repository.GetImageJob(id)
	if err != nil {
		slog.Error(fmt.Sprintf("Error fetching job: %v", err))
		return
	}
	return
}

func (s *service) GetJobsByStatus(status string) (imageJobsResponse []dto.ImageJob, err error) {
	imageJobsResponse, err = s.repository.GetImageJobsByStatus(status)

	if err != nil {
		slog.Error(fmt.Sprintf("Error fetching job: %v", err))
		return
	}
	return
}

func (s *service) RetryJob(id int64) (err error) {
	job, err := s.repository.GetImageJob(id)
	if err != nil {
		slog.Error(fmt.Sprintf("Error fetching job: %v", err))
		return
	}

	if job.Status != "failed" {
		slog.Error("Only failed jobs can be retried")
		return errors.New("only failed jobs can be retried")
	}

	err = s.repository.UpdateJobStatus(id, "pending")
	if err != nil {
		slog.Error(fmt.Sprintf("Error updating job status: %v", err))
		return err
	}

	err = s.rabbitmq.PublishJob(id, job.Filename)
	if err != nil {
		// If publishing fails, revert status to failed
		err = s.repository.UpdateJobStatus(id, "failed")
		if err != nil {
			slog.Error(fmt.Sprintf("Error updating job status: %v", err))
			return err
		}
		slog.Error(fmt.Sprintf("Error publishing job to queue: %v", err))
		return err
	}

	return
}

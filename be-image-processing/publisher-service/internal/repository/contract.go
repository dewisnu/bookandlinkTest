package repository

import (
	"database/sql"
	"publisher-service/pkg/dto"
)

type Repository interface {
	CreateImageJob(filename string, originalSize int64) (int64, error)
	UpdateJobStatus(id int64, status string) error
	GetImageJobs() ([]dto.ImageJob, error)
	GetImageJob(id int64) (dto.ImageJob, error)
	GetImageJobsByStatus(status string) ([]dto.ImageJob, error)
}

type repository struct {
	db   *sql.DB
	conf *repositoryConfig
}

type repositoryConfig struct {
}

type NewRepositoryParams struct {
	Database *sql.DB
}

func NewRepository(params *NewRepositoryParams) Repository {
	if params.Database == nil {
		panic("NewRepository: Database is nil")
	}

	return &repository{
		conf: &repositoryConfig{},
		db:   params.Database,
	}
}

package service

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"publisher-service/internal/config"
	"publisher-service/internal/repository"
	"publisher-service/pkg/dto"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
	HandleUpload(g *gin.Context, files []*multipart.FileHeader) (imageResponse dto.ImageResponse, err error)
	GetJobs() (imageJobsResponse []dto.ImageJob, err error)
	GetJob(id int64) (imageJobResponse dto.ImageJob, err error)
	GetJobsByStatus(status string) (imageJobsResponse []dto.ImageJob, err error)
	RetryJob(id int64) (err error)
	ServeImageUploaded(filename string) (imagePath string, isExist bool, err error)
	ServeImageCompressed(filename string) (imagePath string, isExist bool, err error)
	CompressedUpload(g *gin.Context, file *multipart.FileHeader) (compressedImageResponse dto.CompressedImageResponse, err error)
}

type service struct {
	conf       *serviceConfig
	repository repository.Repository
	rabbitmq   *config.RabbitMQ
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Repository repository.Repository
	RabbitMQ   *config.RabbitMQ
}

func NewService(params *NewServiceParams) Service {
	return &service{
		conf:       &serviceConfig{},
		repository: params.Repository,
		rabbitmq:   params.RabbitMQ,
	}
}

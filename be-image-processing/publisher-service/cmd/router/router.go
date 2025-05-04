package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"publisher-service/cmd/handler"
	"publisher-service/internal/config"
	"publisher-service/internal/service"
)

const (
	HeaderOrigin      = "Origin"
	HeaderContentType = "Content-Type"
	HeaderAccept      = "Accept"
)

type InitRouterParams struct {
	Service service.Service
	Gn      *gin.Engine
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	params.Gn.Use(cors.New(cors.Config{
		AllowOrigins: params.Conf.CorsAllowOrigins,
		AllowHeaders: []string{HeaderOrigin, HeaderContentType, HeaderAccept},
	}))

	params.Gn.GET("/ping", handler.HandlePing(params.Service.Ping))
	params.Gn.POST("/upload", handler.HandleImageUpload(params.Service.HandleUpload))
	params.Gn.GET("/jobs", handler.HandleGetJobs(params.Service.GetJobs))
	params.Gn.GET("jobs/:id", handler.HandleGetJob(params.Service.GetJob))
	params.Gn.GET("/jobs/status/:status", handler.HandleGetJobByStatus(params.Service.GetJobsByStatus))
	params.Gn.POST("/jobs/:id/retry", handler.HandleRetryJobs(params.Service.RetryJob))
	params.Gn.GET("/images-uploaded/:filename", handler.HandleServeImageUploaded(params.Service.ServeImageUploaded))
	params.Gn.GET("/images-compressed/:filename", handler.HandleServeImageCompressed(params.Service.ServeImageCompressed))
	params.Gn.POST("/compressed", handler.HandleCompressedUpload(params.Service.CompressedUpload))
}

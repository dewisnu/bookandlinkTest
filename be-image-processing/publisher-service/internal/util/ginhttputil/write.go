package ginhttputil

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"publisher-service/pkg/dto"
)

func WriteSuccessResponse(gin *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "success"
	}

	gin.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func WriteErrorResponse(gin *gin.Context, status int, err error) {
	gin.JSON(status, dto.BaseResponse{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	})
}

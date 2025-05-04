package handler

import (
	"github.com/gin-gonic/gin"
	"publisher-service/internal/util/ginhttputil"
	"publisher-service/pkg/dto"
)

type PingHandler func() (pingResponse dto.PublicPingResponse)

func HandlePing(handler PingHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		resp := handler()
		ginhttputil.WriteSuccessResponse(g, resp, "")
	}
}

package service

import (
	"publisher-service/pkg/dto"
	"time"
)

func (s *service) Ping() (pingResponse dto.PublicPingResponse) {
	return dto.PublicPingResponse{
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	}
}

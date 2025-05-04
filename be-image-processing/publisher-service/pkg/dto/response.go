package dto

type ErrorResponse struct {
	Success bool
	Message string
}

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

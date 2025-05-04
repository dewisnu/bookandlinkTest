package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"publisher-service/internal/util/ginhttputil"
	"publisher-service/pkg/dto"
	"strconv"
)

type GetJobsHandler func() (imageJobsResponse []dto.ImageJob, err error)
type GetJobHandler func(id int64) (imageJobResponse dto.ImageJob, err error)
type GetJobByStatusHandler func(status string) (imageJobResponse []dto.ImageJob, err error)
type GetJobRetryHandler func(id int64) (err error)

func HandleGetJobs(handler GetJobsHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		resp, err := handler()

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)

		}

		ginhttputil.WriteSuccessResponse(g, resp, "success get jobs")

	}
}

func HandleGetJob(handler GetJobHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		idStr := g.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("invalid job ID"))
			return
		}

		resp, err := handler(id)

		if resp == (dto.ImageJob{}) {
			ginhttputil.WriteErrorResponse(g, http.StatusNotFound, errors.New("job not found"))
			return
		}

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)
			return
		}

		ginhttputil.WriteSuccessResponse(g, resp, "success get job")

	}
}

func HandleGetJobByStatus(handler GetJobByStatusHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		status := g.Param("status")

		if status != "pending" && status != "processing" && status != "completed" && status != "failed" {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("invalid status. Must be one of: pending, processing, complete, failed"))
			return
		}

		resp, err := handler(status)

		if len(resp) == 0 {
			ginhttputil.WriteErrorResponse(g, http.StatusNotFound, errors.New("job not found"))
			return
		}

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)
			return
		}

		ginhttputil.WriteSuccessResponse(g, resp, "success get jobs by status")

	}
}

func HandleRetryJobs(handler GetJobRetryHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		idStr := g.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("invalid job ID"))
			return
		}

		err = handler(id)

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)
			return
		}

		ginhttputil.WriteSuccessResponse(g, nil, "Job queued for retry")

	}
}

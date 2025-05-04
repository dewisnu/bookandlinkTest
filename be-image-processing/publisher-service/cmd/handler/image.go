package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"publisher-service/internal/util/ginhttputil"
	"publisher-service/internal/util/helper"
	"publisher-service/pkg/dto"
)

type ImageUploadHandler func(g *gin.Context, files []*multipart.FileHeader) (imageResponse dto.ImageResponse, err error)
type ServeImageUploadedHandler func(filename string) (imagePath string, isExist bool, err error)
type ServeImageCompressedHandler func(filename string) (imagePath string, isExist bool, err error)
type CompressedUploadHandler func(g *gin.Context, files *multipart.FileHeader) (compressedImageResponse dto.CompressedImageResponse, err error)

func HandleImageUpload(handler ImageUploadHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		form, err := g.MultipartForm()

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("invalid form data"))
			return
		}

		files := form.File["images"]
		if len(files) == 0 {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("no images uploaded"))
			return
		}

		resp, err := handler(g, files)
		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)
			return
		}

		if len(resp.JobIDs) == 0 {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, errors.New("failed to process any uploaded images"))
			return
		}

		ginhttputil.WriteSuccessResponse(g, resp, "success uploaded images")
	}
}

func HandleServeImageUploaded(handler ServeImageUploadedHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		filename := g.Param("filename")

		// Validate filename to avoid path traversal
		if helper.ContainsPathTraversal(filename) {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("invalid filename"))
			return
		}

		resp, isExist, err := handler(filename)

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)
			return
		}

		if isExist {
			ginhttputil.WriteErrorResponse(g, http.StatusNotFound, errors.New("image not found"))
			return
		}

		g.File(resp)
	}
}

func HandleServeImageCompressed(handler ServeImageCompressedHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		filename := g.Param("filename")

		// Validate filename to avoid path traversal
		if helper.ContainsPathTraversal(filename) {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("invalid filename"))
			return
		}

		resp, isExist, err := handler(filename)

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)
			return
		}

		if isExist {
			ginhttputil.WriteErrorResponse(g, http.StatusNotFound, errors.New("image not found"))
			return
		}

		g.File(resp)
	}
}

func HandleCompressedUpload(handler CompressedUploadHandler) gin.HandlerFunc {
	return func(g *gin.Context) {
		file, err := g.FormFile("file")

		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("missing compressed image"))
			return
		}

		if !helper.IsImage(file.Filename) {
			ginhttputil.WriteErrorResponse(g, http.StatusBadRequest, errors.New("uploaded file is not a valid image"))
			return
		}

		resp, err := handler(g, file)
		if err != nil {
			ginhttputil.WriteErrorResponse(g, http.StatusInternalServerError, err)
			return
		}

		ginhttputil.WriteSuccessResponse(g, resp, "success uploaded compressed images")
	}
}

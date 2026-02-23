package handler

import (
	"net/http"

	"github.com/amirullazmi0/kratify-backend/internal/usecase"
	"github.com/amirullazmi0/kratify-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	usecase usecase.AttachmentUsecase
}

func NewAttachmentHandler(usecase usecase.AttachmentUsecase) *AttachmentHandler {
	return &AttachmentHandler{usecase: usecase}
}

// UploadImage godoc
// @Summary Upload an image
// @Description Upload an image file to ImageKit
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} response.Response{data=dto.AttachmentResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/attachments/image [post]
func (h *AttachmentHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to get file from request", err.Error())
		return
	}
	defer file.Close()

	resp, err := h.usecase.UploadImage(file, header.Filename)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload image", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Image uploaded successfully", resp)
}

// UploadDocument godoc
// @Summary Upload a document
// @Description Upload a document/PDF file to ImageKit
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file"
// @Success 200 {object} response.Response{data=dto.AttachmentResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/attachments/document [post]
func (h *AttachmentHandler) UploadDocument(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to get file from request", err.Error())
		return
	}
	defer file.Close()

	resp, err := h.usecase.UploadDocument(file, header.Filename)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload document", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Document uploaded successfully", resp)
}

// UploadProductImage godoc
// @Summary Upload a product image
// @Description Upload a product image file to ImageKit (products folder)
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Product Image file"
// @Success 200 {object} response.Response{data=dto.AttachmentResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/attachments/product-image [post]
func (h *AttachmentHandler) UploadProductImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to get file from request", err.Error())
		return
	}
	defer file.Close()

	resp, err := h.usecase.UploadProductImage(file, header.Filename)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload product image", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product image uploaded successfully", resp)
}

// UploadProfileImage godoc
// @Summary Upload a profile image
// @Description Upload a profile image file to ImageKit (profiles folder)
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Profile Image file"
// @Success 200 {object} response.Response{data=dto.AttachmentResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/attachments/profile-image [post]
func (h *AttachmentHandler) UploadProfileImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to get file from request", err.Error())
		return
	}
	defer file.Close()

	userID := c.GetString("user_id")

	resp, err := h.usecase.UploadProfileImage(file, header.Filename, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload profile image", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile image uploaded successfully", resp)
}

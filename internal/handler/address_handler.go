package handler

import (
	"net/http"

	"github.com/amirullazmi0/kratify-backend/internal/dto"
	"github.com/amirullazmi0/kratify-backend/internal/usecase"
	"github.com/amirullazmi0/kratify-backend/pkg/response"
	"github.com/amirullazmi0/kratify-backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	AddressUsecase usecase.AddressUsecase
}

func NewAddressHandler(addressUsecase usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{AddressUsecase: addressUsecase}
}

func (h *AddressHandler) GetAddressByAuth(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.AddressUsecase.GetAddressByAuth(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Address retrieved successfully", result)
}

func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	addressID := c.Param("address_id")

	result, err := h.AddressUsecase.GetAddressById(addressID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Address retrieved successfully", result)
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	result, err := h.AddressUsecase.CreateAddress(userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Address created successfully", result)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	result, err := h.AddressUsecase.UpdateAddress(userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Address updated successfully", result)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("address_id")

	if err := h.AddressUsecase.DeleteAddress(userID, addressID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Address deleted successfully", nil)
}

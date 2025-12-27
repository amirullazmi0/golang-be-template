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

// GetAddressByAuth godoc
// @Summary Get user addresses
// @Description Get all addresses for authenticated user
// @Tags addresses
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]dto.AddressResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/addresses [get]
func (h *AddressHandler) GetAddressByAuth(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.AddressUsecase.GetAddressByAuth(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Address retrieved successfully", result)
}

// GetAddressByID godoc
// @Summary Get address by ID
// @Description Get specific address by ID
// @Tags addresses
// @Produce json
// @Security BearerAuth
// @Param address_id path string true "Address ID"
// @Success 200 {object} response.Response{data=dto.AddressResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/addresses/{address_id} [get]
func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	addressID := c.Param("address_id")

	result, err := h.AddressUsecase.GetAddressById(addressID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Address retrieved successfully", result)
}

// CreateAddress godoc
// @Summary Create new address
// @Description Create a new address for authenticated user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAddressRequest true "Create Address Request"
// @Success 201 {object} response.Response{data=dto.AddressResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/addresses [post]
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

// UpdateAddress godoc
// @Summary Update address
// @Description Update existing address for authenticated user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateAddressRequest true "Update Address Request"
// @Success 200 {object} response.Response{data=dto.AddressResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/addresses [put]
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

// DeleteAddress godoc
// @Summary Delete address
// @Description Delete address by ID
// @Tags addresses
// @Produce json
// @Security BearerAuth
// @Param address_id path string true "Address ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/addresses/{address_id} [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("address_id")

	if err := h.AddressUsecase.DeleteAddress(userID, addressID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Address deleted successfully", nil)
}

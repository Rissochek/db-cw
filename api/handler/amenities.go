package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/labstack/echo/v4"
)

// @Summary Создать удобство
// @Tags amenities
// @Accept json
// @Produce json
// @Param amenity body AmenityCreate true "Данные удобства"
// @Success 201 {object} model.Amenity
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/amenities [post]
func (h *Handler) CreateAmenity(c echo.Context) error {
	var amenityCreate AmenityCreate
	if err := c.Bind(&amenityCreate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	amenity := model.Amenity{
		Name: amenityCreate.Name,
	}

	if err := h.service.CreateAmenity(c.Request().Context(), &amenity); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, amenity)
}

// @Summary Получить удобство по ID
// @Tags amenities
// @Produce json
// @Param id path int true "Amenity ID"
// @Success 200 {object} model.Amenity
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/amenities/{id} [get]
func (h *Handler) GetAmenityByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid amenity id",
		})
	}

	amenity, err := h.service.GetAmenityByID(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, ErrorNotFound{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, amenity)
}

// @Summary Получить все удобства
// @Tags amenities
// @Produce json
// @Success 200 {array} model.Amenity
// @Failure 500 {object} ErrorInternal
// @Router /api/amenities [get]
func (h *Handler) GetAllAmenities(c echo.Context) error {
	amenities, err := h.service.GetAllAmenities(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, amenities)
}

// @Summary Обновить удобство
// @Tags amenities
// @Accept json
// @Produce json
// @Param id path int true "Amenity ID"
// @Param amenity body AmenityUpdate true "Данные удобства"
// @Success 200 {object} model.Amenity
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/amenities/{id} [put]
func (h *Handler) UpdateAmenity(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid amenity id",
		})
	}

	var amenityUpdate AmenityUpdate
	if err := c.Bind(&amenityUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	amenity := model.Amenity{
		ID:   id,
		Name: amenityUpdate.Name,
	}

	if err := h.service.UpdateAmenity(c.Request().Context(), &amenity); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, ErrorNotFound{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, amenity)
}

// @Summary Удалить удобство
// @Tags amenities
// @Produce json
// @Param id path int true "Amenity ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/amenities/{id} [delete]
func (h *Handler) DeleteAmenity(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid amenity id",
		})
	}

	if err := h.service.DeleteAmenity(c.Request().Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, ErrorNotFound{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, StatusOK{
		Message: "deleted successfully",
	})
}

// @Summary Добавить удобство к объявлению
// @Tags amenities
// @Accept json
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Param amenity_id path int true "Amenity ID"
// @Success 201 {object} model.ListingAmenity
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/listings/{listing_id}/amenities/{amenity_id} [post]
func (h *Handler) AddAmenityToListing(c echo.Context) error {
	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid listing id",
		})
	}

	amenityIDStr := c.Param("amenity_id")
	amenityID, err := strconv.Atoi(amenityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid amenity id",
		})
	}

	if err := h.service.AddAmenityToListing(c.Request().Context(), listingID, amenityID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, ErrorNotFound{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, model.ListingAmenity{
		ListingID: listingID,
		AmenityID: amenityID,
	})
}

// @Summary Удалить удобство из объявления
// @Tags amenities
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Param amenity_id path int true "Amenity ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/listings/{listing_id}/amenities/{amenity_id} [delete]
func (h *Handler) RemoveAmenityFromListing(c echo.Context) error {
	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid listing id",
		})
	}

	amenityIDStr := c.Param("amenity_id")
	amenityID, err := strconv.Atoi(amenityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid amenity id",
		})
	}

	if err := h.service.RemoveAmenityFromListing(c.Request().Context(), listingID, amenityID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, ErrorNotFound{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, StatusOK{
		Message: "deleted successfully",
	})
}

// @Summary Получить удобства объявления
// @Tags amenities
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Success 200 {array} model.Amenity
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/listings/{listing_id}/amenities [get]
func (h *Handler) GetAmenitiesByListingID(c echo.Context) error {
	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid listing id",
		})
	}

	amenities, err := h.service.GetAmenitiesByListingID(c.Request().Context(), listingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, amenities)
}

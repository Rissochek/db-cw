package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/labstack/echo/v4"
)

// @Summary Создать объявление
// @Tags listings
// @Accept json
// @Produce json
// @Param listing body ListingCreate true "Данные объявления"
// @Success 201 {object} model.Listing
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/listings [post]
func (h *Handler) CreateListing(c echo.Context) error {
	var listing model.Listing
	if err := c.Bind(&listing); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := h.service.CreateListing(c.Request().Context(), &listing); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, listing)
}

// @Summary Получить объявление по ID
// @Tags listings
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} model.Listing
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/listings/{id} [get]
func (h *Handler) GetListingByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid listing id",
		})
	}

	listing, err := h.service.GetListingByID(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, listing)
}

// @Summary Обновить объявление
// @Tags listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param listing body ListingUpdate true "Данные объявления"
// @Success 200 {object} model.Listing
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/listings/{id} [put]
func (h *Handler) UpdateListing(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid listing id",
		})
	}

	var listing model.Listing
	if err := c.Bind(&listing); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	listing.ID = id

	if err := h.service.UpdateListing(c.Request().Context(), &listing); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, listing)
}

// @Summary Удалить объявление
// @Tags listings
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/listings/{id} [delete]
func (h *Handler) DeleteListing(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid listing id",
		})
	}

	if err := h.service.DeleteListing(c.Request().Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "listing deleted successfully",
	})
}

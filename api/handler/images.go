package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/labstack/echo/v4"
)

// @Summary Создать изображение
// @Tags images
// @Accept json
// @Produce json
// @Param image body ImageCreate true "Данные изображения"
// @Success 201 {object} model.Image
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/images [post]
func (h *Handler) CreateImage(c echo.Context) error {
	var imageCreate ImageCreate
	if err := c.Bind(&imageCreate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	image := model.Image{
		ListingID:  imageCreate.ListingID,
		ImageURL:   imageCreate.ImageURL,
		IsPrimary:  imageCreate.IsPrimary,
		OrderIndex: imageCreate.OrderIndex,
	}

	if err := h.service.CreateImage(c.Request().Context(), &image); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, image)
}

// @Summary Получить изображение по ID
// @Tags images
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} model.Image
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/images/{id} [get]
func (h *Handler) GetImageByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid image id",
		})
	}

	image, err := h.service.GetImageByID(c.Request().Context(), id)
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

	return c.JSON(http.StatusOK, image)
}

// @Summary Получить изображения по listing ID
// @Tags images
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Success 200 {array} model.Image
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/listings/{listing_id}/images [get]
func (h *Handler) GetImagesByListingID(c echo.Context) error {
	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid listing id",
		})
	}

	images, err := h.service.GetImagesByListingID(c.Request().Context(), listingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, images)
}

// @Summary Обновить изображение
// @Tags images
// @Accept json
// @Produce json
// @Param id path int true "Image ID"
// @Param image body ImageUpdate true "Данные изображения"
// @Success 200 {object} model.Image
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/images/{id} [put]
func (h *Handler) UpdateImage(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid image id",
		})
	}

	var imageUpdate ImageUpdate
	if err := c.Bind(&imageUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	image := model.Image{
		ImageID:    id,
		ImageURL:   imageUpdate.ImageURL,
		IsPrimary:  imageUpdate.IsPrimary,
		OrderIndex: imageUpdate.OrderIndex,
	}

	if err := h.service.UpdateImage(c.Request().Context(), &image); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, ErrorNotFound{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, image)
}

// @Summary Удалить изображение
// @Tags images
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/images/{id} [delete]
func (h *Handler) DeleteImage(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid image id",
		})
	}

	if err := h.service.DeleteImage(c.Request().Context(), id); err != nil {
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

package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// @Summary Создать бронирование
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body BookingCreate true "Данные бронирования"
// @Success 201 {object} model.Booking
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/bookings [post]
func (h *Handler) CreateBooking(c echo.Context) error {
	var booking model.Booking
	if err := c.Bind(&booking); err != nil {
		zap.S().Errorf(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := h.service.CreateBooking(c.Request().Context(), &booking); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, booking)
}

// @Summary Получить бронирование по ID
// @Tags bookings
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} model.Booking
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/bookings/{id} [get]
func (h *Handler) GetBookingByID(c echo.Context) error {
	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid booking id",
		})
	}

	booking, err := h.service.GetBookingByID(c.Request().Context(), bookingID)
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

	return c.JSON(http.StatusOK, booking)
}

// @Summary Обновить бронирование
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param booking body BookingUpdate true "Данные бронирования"
// @Success 200 {object} model.Booking
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/bookings/{id} [put]
func (h *Handler) UpdateBooking(c echo.Context) error {
	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid booking id",
		})
	}

	var booking model.Booking
	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	booking.BookingID = bookingID

	if err := h.service.UpdateBooking(c.Request().Context(), &booking); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, booking)
}

// @Summary Удалить бронирование
// @Tags bookings
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/bookings/{id} [delete]
func (h *Handler) DeleteBooking(c echo.Context) error {
	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid booking id",
		})
	}

	if err := h.service.DeleteBooking(c.Request().Context(), bookingID); err != nil {
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
		"message": "booking deleted successfully",
	})
}

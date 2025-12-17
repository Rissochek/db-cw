package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/labstack/echo/v4"
)

// @Summary Создать платеж
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body PaymentCreate true "Данные платежа"
// @Success 201 {object} model.Payment
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/payments [post]
func (h *Handler) CreatePayment(c echo.Context) error {
	var paymentCreate PaymentCreate
	if err := c.Bind(&paymentCreate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	payment := model.Payment{
		BookingID:     paymentCreate.BookingID,
		PaymentMethod: paymentCreate.PaymentMethod,
		PaymentStatus: paymentCreate.PaymentStatus,
	}

	if err := h.service.CreatePayment(c.Request().Context(), &payment); err != nil {
		if strings.Contains(err.Error(), "must be greater than zero") ||
			strings.Contains(err.Error(), "cannot exceed") ||
			strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, payment)
}

// @Summary Получить платеж по ID
// @Tags payments
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} model.Payment
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/payments/{id} [get]
func (h *Handler) GetPaymentByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid payment id",
		})
	}

	payment, err := h.service.GetPaymentByID(c.Request().Context(), id)
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

	return c.JSON(http.StatusOK, payment)
}

// @Summary Получить платежи по booking ID
// @Tags payments
// @Produce json
// @Param booking_id path int true "Booking ID"
// @Success 200 {array} model.Payment
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/bookings/{booking_id}/payments [get]
func (h *Handler) GetPaymentsByBookingID(c echo.Context) error {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid booking id",
		})
	}

	payments, err := h.service.GetPaymentsByBookingID(c.Request().Context(), bookingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, payments)
}

// @Summary Обновить платеж
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Param payment body PaymentUpdate true "Данные платежа"
// @Success 200 {object} model.Payment
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/payments/{id} [put]
func (h *Handler) UpdatePayment(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid payment id",
		})
	}

	var paymentUpdate PaymentUpdate
	if err := c.Bind(&paymentUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	payment := model.Payment{
		PaymentID:     id,
		PaymentMethod: paymentUpdate.PaymentMethod,
		PaymentStatus: paymentUpdate.PaymentStatus,
	}

	if err := h.service.UpdatePayment(c.Request().Context(), &payment); err != nil {
		if strings.Contains(err.Error(), "must be greater than zero") ||
			strings.Contains(err.Error(), "cannot exceed") ||
			strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, payment)
}

// @Summary Удалить платеж
// @Tags payments
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/payments/{id} [delete]
func (h *Handler) DeletePayment(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid payment id",
		})
	}

	if err := h.service.DeletePayment(c.Request().Context(), id); err != nil {
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

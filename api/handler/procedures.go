package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// @Summary Создать бронирование с платежом через процедуру
// @Tags procedures
// @Accept json
// @Produce json
// @Param request body BookingWithPaymentCreate true "Данные бронирования"
// @Success 201 {object} model.CreateBookingWithPaymentResult
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/procedures/create-booking-with-payment [post]
func (h *Handler) CreateBookingWithPayment(c echo.Context) error {
	var req BookingWithPaymentCreate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	inDate, err := time.Parse(time.RFC3339, req.InDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid in_date format, use RFC3339",
		})
	}

	outDate, err := time.Parse(time.RFC3339, req.OutDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid out_date format, use RFC3339",
		})
	}

	result, err := h.service.CreateBookingWithPayment(
		c.Request().Context(),
		req.ListingID,
		req.GuestID,
		inDate,
		outDate,
		req.PaymentMethod,
	)
	if err != nil {
		if strings.Contains(err.Error(), "listing not found") {
			return c.JSON(http.StatusNotFound, ErrorNotFound{
				Error: "listing not found",
			})
		}
		if strings.Contains(err.Error(), "overlap") || strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "must be later") {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, result)
}

// @Summary Подтвердить платеж через процедуру
// @Tags procedures
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Param request body PaymentConfirmRequest true "Данные подтверждения"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/procedures/payments/{id}/confirm [post]
func (h *Handler) ConfirmPayment(c echo.Context) error {
	idStr := c.Param("id")
	paymentID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid payment id",
		})
	}

	var req PaymentConfirmRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	var transactionID *string
	if req.TransactionID != "" {
		transactionID = &req.TransactionID
	}

	if err := h.service.ConfirmPayment(c.Request().Context(), paymentID, transactionID); err != nil {
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
		Message: "payment confirmed successfully",
	})
}

// @Summary Отменить бронирование с возвратом через процедуру
// @Tags procedures
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/procedures/bookings/{id}/cancel-with-refund [post]
func (h *Handler) CancelBookingWithRefund(c echo.Context) error {
	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid booking id",
		})
	}

	if err := h.service.CancelBookingWithRefund(c.Request().Context(), bookingID); err != nil {
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
		Message: "booking cancelled with refund successfully",
	})
}

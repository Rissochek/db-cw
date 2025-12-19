package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// @Summary Получить общую выручку хоста
// @Tags functions
// @Produce json
// @Param host_id path int true "Host ID"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/functions/hosts/{host_id}/revenue [get]
func (h *Handler) GetHostTotalRevenue(c echo.Context) error {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid host id",
		})
	}

	revenue, err := h.service.GetHostTotalRevenue(c.Request().Context(), hostID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]float64{"total_revenue": revenue})
}

// @Summary Получить общую сумму потраченную гостем
// @Tags functions
// @Produce json
// @Param guest_id path int true "Guest ID"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/functions/guests/{guest_id}/total-spent [get]
func (h *Handler) GetGuestTotalSpent(c echo.Context) error {
	guestIDStr := c.Param("guest_id")
	guestID, err := strconv.Atoi(guestIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid guest id",
		})
	}

	spent, err := h.service.GetGuestTotalSpent(c.Request().Context(), guestID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]float64{"total_spent": spent})
}

// @Summary Получить средний рейтинг хоста
// @Tags functions
// @Produce json
// @Param host_id path int true "Host ID"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/functions/hosts/{host_id}/average-rating [get]
func (h *Handler) GetHostAverageRating(c echo.Context) error {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid host id",
		})
	}

	rating, err := h.service.GetHostAverageRating(c.Request().Context(), hostID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]float64{"average_rating": rating})
}

// @Summary Получить количество активных бронирований объявления
// @Tags functions
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Success 200 {object} map[string]int
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/functions/listings/{listing_id}/active-bookings [get]
func (h *Handler) GetListingActiveBookingsCount(c echo.Context) error {
	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid listing id",
		})
	}

	count, err := h.service.GetListingActiveBookingsCount(c.Request().Context(), listingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]int{"active_bookings_count": count})
}

// @Summary Получить статистический отчет по объявлениям
// @Tags functions
// @Produce json
// @Success 200 {array} model.ListingStatisticsReport
// @Failure 500 {object} ErrorInternal
// @Router /api/reports/listings-statistics [get]
func (h *Handler) GetListingsStatisticsReport(c echo.Context) error {
	reports, err := h.service.GetListingsStatisticsReport(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reports)
}

// @Summary Получить отчет о производительности хостов
// @Tags functions
// @Produce json
// @Success 200 {array} model.HostPerformanceReport
// @Failure 500 {object} ErrorInternal
// @Router /api/reports/hosts-performance [get]
func (h *Handler) GetHostsPerformanceReport(c echo.Context) error {
	reports, err := h.service.GetHostsPerformanceReport(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reports)
}

// @Summary Получить отчет по бронированиям
// @Tags functions
// @Produce json
// @Param start_date query string false "Start Date" format(date-time) example(2025-01-01T00:00:00Z)
// @Param end_date query string false "End Date" format(date-time) example(2025-12-31T23:59:59Z)
// @Success 200 {array} model.BookingReport
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/reports/bookings [get]
func (h *Handler) GetBookingsReport(c echo.Context) error {
	var startDate, endDate *time.Time

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: "invalid start_date format, use RFC3339",
			})
		}
		startDate = &parsed
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: "invalid end_date format, use RFC3339",
			})
		}
		endDate = &parsed
	}

	reports, err := h.service.GetBookingsReport(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reports)
}

// @Summary Получить сводный отчет по платежам
// @Tags functions
// @Produce json
// @Param start_date query string false "Start Date" format(date-time) example(2025-01-01T00:00:00Z)
// @Param end_date query string false "End Date" format(date-time) example(2025-12-31T23:59:59Z)
// @Success 200 {array} model.PaymentSummaryReport
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/reports/payments-summary [get]
func (h *Handler) GetPaymentsSummaryReport(c echo.Context) error {
	var startDate, endDate *time.Time

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: "invalid start_date format, use RFC3339",
			})
		}
		startDate = &parsed
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: "invalid end_date format, use RFC3339",
			})
		}
		endDate = &parsed
	}

	reports, err := h.service.GetPaymentsSummaryReport(c.Request().Context(), startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reports)
}

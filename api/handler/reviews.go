package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/labstack/echo/v4"
)

// @Summary Создать отзыв
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body ReviewCreate true "Данные отзыва"
// @Success 201 {object} model.Review
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/reviews [post]
func (h *Handler) CreateReview(c echo.Context) error {
	var review model.Review
	if err := c.Bind(&review); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := h.service.CreateReview(c.Request().Context(), &review); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, review)
}

// @Summary Получить отзыв по ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} model.Review
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/reviews/{id} [get]
func (h *Handler) GetReviewByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid review id",
		})
	}

	review, err := h.service.GetReviewByID(c.Request().Context(), id)
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

	return c.JSON(http.StatusOK, review)
}

// @Summary Обновить отзыв
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body ReviewUpdate true "Данные отзыва"
// @Success 200 {object} model.Review
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/reviews/{id} [put]
func (h *Handler) UpdateReview(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid review id",
		})
	}

	var review model.Review
	if err := c.Bind(&review); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	review.ID = id

	if err := h.service.UpdateReview(c.Request().Context(), &review); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, review)
}

// @Summary Удалить отзыв
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/reviews/{id} [delete]
func (h *Handler) DeleteReview(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid review id",
		})
	}

	if err := h.service.DeleteReview(c.Request().Context(), id); err != nil {
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
		"message": "review deleted successfully",
	})
}

// @Summary Batch импорт отзывов
// @Tags reviews
// @Accept json
// @Produce json
// @Param reviews body []ReviewCreate true "Массив отзывов"
// @Success 201 {object} map[string]int "Количество созданных отзывов"
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/reviews/batch [post]
func (h *Handler) BatchImportReviews(c echo.Context) error {
	var reviewsCreate []ReviewCreate
	if err := c.Bind(&reviewsCreate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	reviews := make([]model.Review, len(reviewsCreate))
	for i := range reviewsCreate {
		reviews[i].BookingID = reviewsCreate[i].BookingID
		// reviews[i].UserID = reviewsCreate[i].UserID
		reviews[i].Score = reviewsCreate[i].Score
		if reviewsCreate[i].Text != "" {
			reviews[i].Text = reviewsCreate[i].Text
		}
	}

	if err := h.service.CreateReviews(c.Request().Context(), reviews); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]int{
		"created": len(reviews),
	})
}

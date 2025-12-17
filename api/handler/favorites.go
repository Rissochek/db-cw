package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/labstack/echo/v4"
)

// @Summary Добавить объявление в избранное
// @Tags favorites
// @Accept json
// @Produce json
// @Param favorite body FavoriteCreate true "Данные избранного"
// @Success 201 {object} model.Favorite
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/favorites [post]
func (h *Handler) CreateFavorite(c echo.Context) error {
	var favoriteCreate FavoriteCreate
	if err := c.Bind(&favoriteCreate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid request body",
		})
	}

	favorite := model.Favorite{
		UserID:    favoriteCreate.UserID,
		ListingID: favoriteCreate.ListingID,
	}

	if err := h.service.CreateFavorite(c.Request().Context(), &favorite); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusBadRequest, ErrorBadRequest{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, favorite)
}

// @Summary Получить избранное по ID
// @Tags favorites
// @Produce json
// @Param id path int true "Favorite ID"
// @Success 200 {object} model.Favorite
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/favorites/{id} [get]
func (h *Handler) GetFavoriteByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid favorite id",
		})
	}

	favorite, err := h.service.GetFavoriteByID(c.Request().Context(), id)
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

	return c.JSON(http.StatusOK, favorite)
}

// @Summary Получить избранное пользователя
// @Tags favorites
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} model.Favorite
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/users/{user_id}/favorites [get]
func (h *Handler) GetFavoritesByUserID(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid user id",
		})
	}

	favorites, err := h.service.GetFavoritesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorInternal{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, favorites)
}

// @Summary Удалить избранное по ID
// @Tags favorites
// @Produce json
// @Param id path int true "Favorite ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/favorites/{id} [delete]
func (h *Handler) DeleteFavorite(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid favorite id",
		})
	}

	if err := h.service.DeleteFavorite(c.Request().Context(), id); err != nil {
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

// @Summary Удалить избранное по user_id и listing_id
// @Tags favorites
// @Produce json
// @Param user_id path int true "User ID"
// @Param listing_id path int true "Listing ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/users/{user_id}/favorites/{listing_id} [delete]
func (h *Handler) DeleteFavoriteByUserAndListing(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid user id",
		})
	}

	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorBadRequest{
			Error: "invalid listing id",
		})
	}

	if err := h.service.DeleteFavoriteByUserAndListing(c.Request().Context(), userID, listingID); err != nil {
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

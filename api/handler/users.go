package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/Rissochek/db-cw/internal/utils"
	"github.com/labstack/echo/v4"
)

// @Summary Создать пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserCreate true "Данные пользователя"
// @Success 201 {object} UserReturn
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	hashedPassword, err := utils.GenerateHash(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to hash password",
		})
	}
	user.Password = hashedPassword

	if err := h.service.CreateUser(c.Request().Context(), &user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// @Summary Получить пользователя по ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserReturn
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/users/{id} [get]
func (h *Handler) GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), id)
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

	return c.JSON(http.StatusOK, user)
}

// @Summary Обновить пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UserUpdate true "Данные пользователя"
// @Success 200 {object} UserReturn
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
	}

	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	user.ID = id

	if user.Password != "" {
		hashedPassword, err := utils.GenerateHash(user.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to hash password",
			})
		}
		user.Password = hashedPassword
	}

	if err := h.service.UpdateUser(c.Request().Context(), &user); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary Удалить пользователя
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} StatusOK
// @Failure 400 {object} ErrorBadRequest
// @Failure 404 {object} ErrorNotFound
// @Failure 500 {object} ErrorInternal
// @Router /api/users/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
	}

	if err := h.service.DeleteUser(c.Request().Context(), id); err != nil {
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
		"message": "user deleted successfully",
	})
}

// @Summary Batch импорт пользователей
// @Tags users
// @Accept json
// @Produce json
// @Param users body []UserCreate true "Массив пользователей"
// @Success 201 {object} map[string]int "Количество созданных пользователей"
// @Failure 400 {object} ErrorBadRequest
// @Failure 500 {object} ErrorInternal
// @Router /api/users/batch [post]
func (h *Handler) BatchImportUsers(c echo.Context) error {
	var usersCreate []UserCreate
	if err := c.Bind(&usersCreate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	users := make([]model.User, len(usersCreate))
	for i := range usersCreate {
		users[i].Email = usersCreate[i].Email
		users[i].FirstName = usersCreate[i].FirstName
		users[i].SecondName = usersCreate[i].SecondName

		hashedPassword, err := utils.GenerateHash(usersCreate[i].Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to hash password",
			})
		}
		users[i].Password = hashedPassword
	}

	if err := h.service.CreateUsers(c.Request().Context(), users); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]int{
		"created": len(users),
	})
}

package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
	"net/http"
)

func (r *Rest) CreateUser(c echo.Context) error {
	user := new(models.User)
	err := c.Bind(user)
	if err != nil {
		res := logs.OutputRestError("bind error", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.CreateUser(c.Request().Context(), user)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, user)
}

func (r *Rest) GetUserByLogin(c echo.Context) error {
	login := c.Param("login")

	user, err := r.Handler.GetUserByLogin(c.Request().Context(), login)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, user)
}

func (r *Rest) SetUserState(c echo.Context) error {
	params := new(proto.UsersFilter)
	err := c.Bind(params)
	if err != nil {
		res := logs.OutputRestError("bind error", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = r.Handler.SetUserState(c.Request().Context(), params)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, "")
}

func (r *Rest) MarkUserAsDeleted(c echo.Context) error {
	userId := c.Param("user_id")

	err := r.Handler.MarkUserAsDeleted(c.Request().Context(), userId)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, "")
}

package handlers_http_internal_users_v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	pkg_http "github.com/subeecore/pkg/http"
	entities_user_v1 "github.com/subeecore/subee-core-svc/internal/entities/user/v1"
)

type CreateUserRequest struct {
	ID       string `json:"$id"`
	Email    string `json:"email"`
	Username string `json:"name"`
}

func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, pkg_http.NewHTTPResponse(http.StatusBadRequest, pkg_http.MessageBadRequestError, nil))
	}

	_, err := h.service.CreateUser(ctx, &entities_user_v1.CreateUserRequest{
		ExternalID: req.ID,
		Email:      req.Email,
		Username:   req.Username,
	})
	if err != nil {
		return c.JSON(pkg_http.TranslateError(ctx, err))
	}

	return c.JSON(http.StatusOK, pkg_http.NewHTTPResponse(http.StatusOK, pkg_http.MessageSuccess, nil))
}

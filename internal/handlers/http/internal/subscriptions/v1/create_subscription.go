package handlers_http_internal_subscriptions_v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	pkg_http "github.com/subeecore/pkg/http"
	entities_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/entities/subscriptions/v1"
)

type CreateSubscriptionResponse struct {
	Subscription *entities_subscriptions_v1.Subscription `json:"subscription"`
}

func (h *Handler) CreateSubscription(c echo.Context) error {
	ctx := c.Request().Context()

	req := &entities_subscriptions_v1.CreateSubscriptionRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, pkg_http.NewHTTPResponse(http.StatusBadRequest, pkg_http.MessageBadRequestError, nil))
	}

	subscription, err := h.service.CreateSubscription(ctx, req)
	if err != nil {
		return c.JSON(pkg_http.TranslateError(ctx, err))
	}

	return c.JSON(http.StatusCreated, pkg_http.NewHTTPResponse(http.StatusCreated, pkg_http.MessageSuccess, &CreateSubscriptionResponse{
		Subscription: subscription,
	}))
}

package handlers_http_internal_subscriptions_v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	pkg_http "github.com/subeecore/pkg/http"
	entities_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/entities/subscriptions/v1"
)

type FetchSubscriptionResponse struct {
	Subscriptions []*entities_subscriptions_v1.Subscription `json:"subscriptions"`
}

func (h *Handler) FetchSubscriptions(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Param("user_id")
	if userID == "" {
		log.Error().Msg("handlers.http.private.subscriptions.v1.fetch_subscriptions.Handler.FetchSubscriptions: can not get user_id from context")
		return c.JSON(http.StatusBadRequest, pkg_http.NewHTTPResponse(http.StatusBadRequest, pkg_http.MessageBadRequestError, nil))
	}

	subscriptions, err := h.service.FetchSubscriptions(ctx, userID)
	if err != nil {
		return c.JSON(pkg_http.TranslateError(ctx, err))
	}

	return c.JSON(http.StatusOK, pkg_http.NewHTTPResponse(http.StatusOK, pkg_http.MessageSuccess, &FetchSubscriptionResponse{
		Subscriptions: subscriptions,
	}))
}

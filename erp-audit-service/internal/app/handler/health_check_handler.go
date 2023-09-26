package handler

import (
	"net/http"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/healthx"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/global"
	"github.com/labstack/echo/v4"
)

type HealthCheckHandler struct {
	Option global.HandlerOptions
}

// URLMapping implements router.RouteHandlers
func (h *HealthCheckHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.GET("", h.Index)
}

func (h HealthCheckHandler) Index(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var healthCheckStatus healthx.HealthCheckStatus
	healthCheckStatus, err = healthx.Check(&h.Option.Common)
	ctx.JSON(http.StatusOK, edenlabs.FormatResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: h.Option.Common.Config.App.Name,
		Data:    healthCheckStatus,
	})
	return
}

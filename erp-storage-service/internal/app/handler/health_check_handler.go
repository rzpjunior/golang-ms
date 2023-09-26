package handler

import (
	"net/http"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/project-version3/erp-services/erp-storage-service/global"
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
	ctx.Message(http.StatusOK, "success", h.Option.Common.Config.App.Name)
	return
}

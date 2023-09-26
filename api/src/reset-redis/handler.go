package reset_redis

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	key := ctx.QueryParam("KEY")
	if key == "reset_database_redis" {
		ctx.ResponseData = dbredis.Redis.DeleteAllCache()
	} else {
		e = echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	return ctx.Serve(e)
}

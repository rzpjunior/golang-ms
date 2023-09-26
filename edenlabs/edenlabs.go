package edenlabs

import (
	"github.com/labstack/echo/v4"
)

// New creates an instance of Echo.
func New() (e *echo.Echo) {
	e = echo.New()
	return
}

// HTTPErrorHandler invokes the default HTTP error handler.
func HTTPErrorHandler(err error, c echo.Context) {
	if !c.Response().Committed {
		ctx, ok := c.(*Context)
		if !ok {
			ctx = NewContext(c)
		}
		ctx.Serve(err)
	}
}

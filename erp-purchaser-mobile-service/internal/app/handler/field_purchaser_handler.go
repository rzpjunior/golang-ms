package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type FieldPurchaserHandler struct {
	Option                global.HandlerOptions
	ServiceFieldPurchaser service.IFieldPurchaserService
	ServicesAuth          service.IAuthService
}

// URLMapping declare endpoint with handler function.
func (h *FieldPurchaserHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceFieldPurchaser = service.NewFieldPurchaserService()
	h.ServicesAuth = service.NewAuthService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
}

func (h FieldPurchaserHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var session dto.UserResponse

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	// siteID := ctx.GetParamInt("siteID")
	siteID := ctx.GetParamString("siteID")
	if session, err = h.ServicesAuth.Session(ctx.Request().Context()); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if session.Site.ID != "" {
		siteID = session.Site.ID
	}
	var fieldPurchasers []*dto.FieldPurchaserResponse
	var total int64
	fieldPurchasers, total, err = h.ServiceFieldPurchaser.Get(ctx.Request().Context(), &dto.FieldPurchaserListRequest{
		Limit:   int32(page.Limit),
		Offset:  int32(page.Offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
		// SiteID:  int32(siteID),
		SiteIDGp: siteID,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(fieldPurchasers, total, page)

	return ctx.Serve(err)
}

func (h FieldPurchaserHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var fieldPurchaser *dto.FieldPurchaserResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	fieldPurchaser, err = h.ServiceFieldPurchaser.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = fieldPurchaser

	return ctx.Serve(err)
}

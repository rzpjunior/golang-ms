package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ConsolidatedShipmentHandler struct {
	Option                      global.HandlerOptions
	ServiceConsolidatedShipment service.IConsolidatedShipmentService
	ServicesAuth                service.IAuthService
}

// URLMapping declare endpoint with handler function.
func (h *ConsolidatedShipmentHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceConsolidatedShipment = service.NewConsolidatedShipmentService()
	h.ServicesAuth = service.NewAuthService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
	r.POST("/signature", h.Signature, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/print/:id", h.Print, cMiddleware.Authorized("purchaser_app"))
	r.POST("", h.Create, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/:id", h.Update, cMiddleware.Authorized("purchaser_app"))
}

func (h ConsolidatedShipmentHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var session dto.UserResponse
	var employeeID int64
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
	siteID := ctx.GetParamString("site_id")
	createdAtFrom := ctx.GetParamDate("created_at_from")
	createdAtTo := ctx.GetParamDate("created_at_to")
	employeeID = utils.ToInt64(ctx.GetParamInt("created_by"))

	var consolidatedShipments []*dto.ConsolidatedShipmentResponse
	var total int64

	if session, err = h.ServicesAuth.Session(ctx.Request().Context()); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if session.Site.ID != "" {
		siteID = session.Site.ID
	}
	if session.MainRole == "Field Purchaser" {
		employeeID = session.ID
	}

	consolidatedShipments, total, err = h.ServiceConsolidatedShipment.Get(ctx.Request().Context(), &dto.ConsolidatedShipmentRequestList{
		Limit:         int32(page.Limit),
		Offset:        int32(page.Offset),
		Status:        int32(status),
		Search:        search,
		OrderBy:       orderBy,
		SiteID:        siteID,
		CreatedAtFrom: createdAtFrom,
		CreatedAtTo:   createdAtTo,
		CreatedBy:     employeeID,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(consolidatedShipments, total, page)

	return ctx.Serve(err)
}

func (h ConsolidatedShipmentHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var consolidatedShipment *dto.ConsolidatedShipmentResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	consolidatedShipment, err = h.ServiceConsolidatedShipment.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = consolidatedShipment

	return ctx.Serve(err)
}

func (h ConsolidatedShipmentHandler) Signature(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var consolidatedShipmentSignature *dto.ConsolidatedShipmentSignatureResponse

	var req dto.ConsolidatedShipmentSignatureRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	consolidatedShipmentSignature, err = h.ServiceConsolidatedShipment.Signature(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = consolidatedShipmentSignature

	return ctx.Serve(err)
}

func (h ConsolidatedShipmentHandler) Print(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var consolidatedShipment *dto.ConsolidatedShipmentResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	consolidatedShipment, err = h.ServiceConsolidatedShipment.Print(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = consolidatedShipment

	return ctx.Serve(err)
}

func (h ConsolidatedShipmentHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var consolidatedShipment *dto.ConsolidatedShipmentResponse

	var req dto.ConsolidatedShipmentRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	consolidatedShipment, err = h.ServiceConsolidatedShipment.Create(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = consolidatedShipment

	return ctx.Serve(err)
}

func (h ConsolidatedShipmentHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var consolidatedShipment *dto.ConsolidatedShipmentResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.ConsolidatedShipmentRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	consolidatedShipment, err = h.ServiceConsolidatedShipment.Update(ctx.Request().Context(), &req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = consolidatedShipment

	return ctx.Serve(err)
}

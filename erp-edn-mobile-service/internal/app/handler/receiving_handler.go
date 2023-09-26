package handler

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/labstack/echo/v4"
)

type ReceivingHandler struct {
	Option            global.HandlerOptions
	ServicesReceiving service.IReceivingService
}

// URLMapping implements router.RouteHandlers
func (h *ReceivingHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesReceiving = service.NewServiceReceiving()

	cMiddleware := middleware.NewMiddleware()
	// r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	// r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	// r.POST("", h.Create, cMiddleware.Authorized("edn_app"))
	// r.PUT("/confirm/:id", h.Confirm, cMiddleware.Authorized("edn_app"))
	// r.POST("", h.Create, cMiddleware.Authorized("edn_app"))
	// r.PUT("/confirm/:id", h.Confirm, cMiddleware.Authorized("edn_app"))

	// gp integrated
	r.GET("", h.GetGP, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGP, cMiddleware.Authorized("edn_app"))
	r.POST("/gp", h.CreateGP, cMiddleware.Authorized("edn_app"))
	r.PUT("/gp", h.UpdateGP, cMiddleware.Authorized("edn_app"))
}

func (h ReceivingHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

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
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var po []*dto.ReceivingResponse
	var total int64
	po, err = h.ServicesReceiving.Get(ctx.Request().Context(), dto.ReceivingListRequest{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(po, total, page)

	return ctx.Serve(err)
}

func (h ReceivingHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ReceivingResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesReceiving.GetById(ctx.Request().Context(), dto.ReceivingDetailRequest{
		Id: id,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ReceivingHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ReceivingResponse

	var req dto.CreateReceivingRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesReceiving.Create(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ReceivingHandler) Confirm(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ReceivingResponse
	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.ConfirmReceivingRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.Id = fmt.Sprintf("%d", id)
	po, err = h.ServicesReceiving.Confirm(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ReceivingHandler) GetGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	poprctnm := ctx.GetParamString("poprctnm")
	doctype := ctx.GetParamString("doctype")

	var po []*dto.ReceivingResponse
	var total int64
	po, total, err = h.ServicesReceiving.GetListGP(ctx.Request().Context(), dto.GetGoodsReceiptGPListRequest{
		Limit:    int32(page.Limit),
		Offset:   int32(page.Offset),
		Poprctnm: poprctnm,
		Doctype:  doctype,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(po, total, page)

	return ctx.Serve(err)
}

func (h ReceivingHandler) DetailGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *bridgeService.GoodsReceiptGP

	var id string
	id = ctx.Param("id")

	po, err = h.ServicesReceiving.GetDetailGP(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ReceivingHandler) CreateGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *bridgeService.CreateTransferRequestGPResponse

	var req dto.CreateGoodsReceiptGPRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesReceiving.CreateGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ReceivingHandler) UpdateGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *bridgeService.CreateTransferRequestGPResponse

	var req dto.UpdateGoodsReceiptGPRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesReceiving.UpdateGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

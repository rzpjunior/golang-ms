package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/labstack/echo/v4"
)

type PurchaseOrderHandler struct {
	Option                global.HandlerOptions
	ServicesPurchaseOrder service.IPurchaseOrderService
	ServiceAuth           service.IAuthService
}

// URLMapping implements router.RouteHandlers
func (h *PurchaseOrderHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPurchaseOrder = service.NewServicePurchaseOrder()
	h.ServiceAuth = service.NewServiceAuth()

	cMiddleware := middleware.NewMiddleware()
	// r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	// r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	// r.POST("", h.Create, cMiddleware.Authorized("edn_app"))
	// r.PUT("/commit/:id", h.Commit, cMiddleware.Authorized("edn_app"))
	// r.PUT("/cancel/:id", h.Cancel, cMiddleware.Authorized("edn_app"))
	// r.PUT("/:id", h.Update, cMiddleware.Authorized("edn_app"))
	// r.PUT("/update-product/:id", h.UpdateProduct, cMiddleware.Authorized("edn_app"))
	r.POST("", h.CreateGP, cMiddleware.Authorized("edn_app"))
	r.PUT("", h.UpdateGP, cMiddleware.Authorized("edn_app"))
	r.GET("", h.GetGP, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGP, cMiddleware.Authorized("edn_app"))
	r.PUT("/commit", h.CommitGP, cMiddleware.Authorized("edn_app"))
}

func (h PurchaseOrderHandler) Get(c echo.Context) (err error) {
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

	var po []*dto.PurchaseOrderResponse
	var total int64
	po, err = h.ServicesPurchaseOrder.Get(ctx.Request().Context(), dto.PurchaseOrderListRequest{
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

func (h PurchaseOrderHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.PurchaseOrderResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesPurchaseOrder.GetById(ctx.Request().Context(), dto.PurchaseOrderDetailRequest{
		Id: id,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.PurchaseOrderResponse

	var req dto.CreatePurchaseOrderRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesPurchaseOrder.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Commit(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesPurchaseOrder.Commit(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = nil

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Cancel(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		id  int64
		req dto.CancelPurchaseOrderRequest
	)
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.Id = id
	err = h.ServicesPurchaseOrder.Cancel(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = nil

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.PurchaseOrderResponse

	var req dto.UpdatePurchaseOrderRequest
	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesPurchaseOrder.Update(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) UpdateProduct(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.PurchaseOrderResponse

	var req dto.UpdateProductPurchaseOrderRequest
	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesPurchaseOrder.UpdateProduct(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) CreateGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *bridgeService.CreatePurchaseOrderGPResponse

	var req dto.CreatePurchaseOrderGPRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesPurchaseOrder.CreateGP(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) UpdateGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *bridgeService.CreatePurchaseOrderGPResponse

	var req dto.CreatePurchaseOrderGPRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesPurchaseOrder.UpdateGP(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) GetGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var (
		page    *edenlabs.Paginator
		po      []*dto.PurchaseOrderResponse
		total   int64
		session dto.UserResponse
	)

	if session, err = h.ServiceAuth.Session(ctx.Request().Context()); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	ponumber := ctx.GetParamString("code")
	ponumberLike := ctx.GetParamString("codelike")
	vendorId := ctx.GetParamString("vendor")
	siteId := ctx.GetParamString("site")
	orderBy := ctx.GetParamString("orderby")
	reqDateFrom := ctx.GetParamDate("reqdatefrom")
	reqDateTo := ctx.GetParamDate("reqdateto")
	status := ctx.GetParamInt("status")

	if siteId == "" {
		siteId = session.SiteID
	}

	po, total, err = h.ServicesPurchaseOrder.GetListGP(ctx.Request().Context(), dto.GetPurchaseOrderGPListRequest{
		Limit:        int32(page.PerPage),
		Offset:       int32(page.Page - 1),
		Ponumber:     ponumber,
		PonumberLike: ponumberLike,
		Vendorid:     vendorId,
		Locncode:     siteId,
		ReqDateFrom:  reqDateFrom,
		ReqDateTo:    reqDateTo,
		OrderBy:      orderBy,
		Postatus:     status,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(po, total, page)

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) DetailGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.PurchaseOrderResponse

	var id string
	id = ctx.Param("id")

	po, err = h.ServicesPurchaseOrder.GetDetailGP(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) CommitGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		req dto.CommitPurchaseOrderGPRequest
		po  *bridgeService.CreateTransferRequestGPResponse
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesPurchaseOrder.CommitPurchaseOrderGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

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

type ItemTransferHandler struct {
	Option               global.HandlerOptions
	ServicesItemTransfer service.IItemTransferService
	ServiceAuth          service.IAuthService
}

// URLMapping implements router.RouteHandlers
func (h *ItemTransferHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesItemTransfer = service.NewServiceItemTransfer()
	h.ServiceAuth = service.NewServiceAuth()

	cMiddleware := middleware.NewMiddleware()
	// r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	// r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	// r.POST("", h.Create, cMiddleware.Authorized("edn_app"))
	// r.PUT("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	// r.PUT("/commit/:id", h.Commit, cMiddleware.Authorized("edn_app"))

	// gp
	r.GET("/request", h.GetTransferRequestGP, cMiddleware.Authorized("edn_app"))
	r.GET("/request/:id", h.TransferRequestDetailGP, cMiddleware.Authorized("edn_app"))

	r.POST("/request", h.CreateTransferRequestGP, cMiddleware.Authorized("edn_app"))
	r.PUT("/request", h.UpdateTransferRequestGP, cMiddleware.Authorized("edn_app"))
	r.PUT("/request/commit", h.CommitTransferRequestGP, cMiddleware.Authorized("edn_app"))

	r.GET("/gp/intransit", h.GetInTransitTransferGP, cMiddleware.Authorized("edn_app"))
	r.GET("/intransit/:id", h.InTransitTransferDetailGP, cMiddleware.Authorized("edn_app"))

	r.PUT("/intransit", h.UpdateInTransitTransferGP, cMiddleware.Authorized("edn_app"))
}

func (h ItemTransferHandler) Get(c echo.Context) (err error) {
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

	var po []*dto.ItemTransferResponse
	var total int64
	po, err = h.ServicesItemTransfer.Get(ctx.Request().Context(), dto.ItemTransferListRequest{
		Limit:   int32(page.Limit),
		Offset:  int32(page.Offset),
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

func (h ItemTransferHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ItemTransferResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesItemTransfer.GetById(ctx.Request().Context(), dto.ItemTransferDetailRequest{
		Id: id,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ItemTransferHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ItemTransferResponse

	var req dto.CreateItemTransferRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesItemTransfer.Create(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ItemTransferHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ItemTransferResponse

	var req dto.UpdateItemTransferRequest
	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesItemTransfer.Update(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ItemTransferHandler) Commit(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ItemTransferResponse

	var req dto.CommitItemTransferRequest
	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesItemTransfer.Commit(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ItemTransferHandler) GetInTransitTransferGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	orddocid := ctx.GetParamString("orddocid")
	ivmTrType := ctx.GetParamString("ivm_tr_type")
	ordrdate := ctx.GetParamString("ordrdate")
	trnsfloc := ctx.GetParamString("trnsfloc")
	locncode := ctx.GetParamString("locncode")
	requestDate := ctx.GetParamString("request_date")
	etadte := ctx.GetParamString("etadte")
	status := ctx.GetParamInt("status")

	var po []*bridgeService.InTransitTransferGP
	var total int64
	po, total, err = h.ServicesItemTransfer.GetInTransitTransferListGP(ctx.Request().Context(), &dto.GetInTransitTransferGPListRequest{
		Limit:       int32(page.Limit),
		Offset:      int32(page.Offset),
		Orddocid:    orddocid,
		IvmTrType:   ivmTrType,
		Ordrdate:    ordrdate,
		Trnsfloc:    trnsfloc,
		Locncode:    locncode,
		RequestDate: requestDate,
		Etadte:      etadte,
		Status:      int32(status),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(po, total, page)

	return ctx.Serve(err)
}

func (h ItemTransferHandler) InTransitTransferDetailGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var itt *dto.ItemTransferResponse

	var id string
	id = ctx.Param("id")

	itt, err = h.ServicesItemTransfer.GetInTransitTransferDetailGP(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = itt

	return ctx.Serve(err)
}

func (h ItemTransferHandler) GetTransferRequestGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var (
		page    *edenlabs.Paginator
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
	// cannot filter like gp
	// trnumberLike := ctx.GetParamString("codelike")
	trnumber := ctx.GetParamString("code")
	siteFromId := ctx.GetParamString("site_from")
	siteToId := ctx.GetParamString("site_to")
	orderBy := ctx.GetParamString("orderby")
	reqDateFrom := ctx.GetParamDate("reqdatefrom")
	reqDateTo := ctx.GetParamDate("reqdateto")
	status := ctx.GetParamInt("status")

	if siteToId == "" {
		siteToId = session.SiteID
	}

	var po []*dto.ItemTransferResponse
	var total int64
	po, total, err = h.ServicesItemTransfer.GetTransferRequestListGP(ctx.Request().Context(), &dto.GetTransferRequestGPListRequest{
		Limit:           int32(page.PerPage),
		Offset:          int32(page.Page - 1),
		Docnumbr:        trnumber,
		RequestDateFrom: reqDateFrom,
		RequestDateTo:   reqDateTo,
		IvmLocncodeFrom: siteFromId,
		IvmLocncodeTo:   siteToId,
		OrderBy:         orderBy,
		Status:          int32(status),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(po, total, page)

	return ctx.Serve(err)
}

func (h ItemTransferHandler) TransferRequestDetailGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.ItemTransferResponse

	var id string
	id = ctx.Param("id")

	po, err = h.ServicesItemTransfer.GetTransferRequestDetailGP(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h ItemTransferHandler) CreateTransferRequestGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		req dto.CreateTransferRequestGPRequest
		tr  *bridgeService.CreateTransferRequestGPResponse
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	tr, err = h.ServicesItemTransfer.CreateTransferRequestGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = tr

	return ctx.Serve(err)
}

func (h ItemTransferHandler) UpdateTransferRequestGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		req dto.UpdateTransferRequestGPRequest
		tr  *bridgeService.CreateTransferRequestGPResponse
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	tr, err = h.ServicesItemTransfer.UpdateTransferRequestGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = tr

	return ctx.Serve(err)
}

func (h ItemTransferHandler) UpdateInTransitTransferGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		req dto.UpdateInTransiteTransferGPRequest
		itt *bridgeService.UpdateInTransitTransferGPResponse
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	itt, err = h.ServicesItemTransfer.UpdateInTransitTransferGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = itt

	return ctx.Serve(err)
}

func (h ItemTransferHandler) CommitTransferRequestGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		req dto.CommitTransferRequestGPRequest
		tr  *bridgeService.CommitTransferRequestGPResponse
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	tr, err = h.ServicesItemTransfer.CommitTransferRequestGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = tr

	return ctx.Serve(err)
}

package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type TransactionHistoryHandler struct {
	Option                     global.HandlerOptions
	ServicesTransactionHistory service.ITransactionHistoryService
}

// URLMapping declare endpoint with handler function.
func (h *TransactionHistoryHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	cMiddleware := middleware.NewMiddleware()
	h.ServicesTransactionHistory = service.NewTransactionHistoryService()

	r.POST("", h.read, cMiddleware.Authorized("private"))
	r.POST("/detail", h.readDetail, cMiddleware.Authorized("private"))
	r.POST("/detail/invoice", h.readDetailSI, cMiddleware.Authorized("private"))
}

func (h TransactionHistoryHandler) read(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.RequestGetHistoryTransaction

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	var page *edenlabs.Paginator
	ctx.ResponseFormat.Page = int(req.Offset) + 1
	ctx.ResponseFormat.PerPage = int(req.Limit)
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	if req.Session, err = service.CustomerSession(ctx); err == nil {
		res, total, e := h.ServicesTransactionHistory.Get(ctx.Request().Context(), req)
		ctx.DataList(res, total, page)
		if e != nil {
			return ctx.Serve(e)
		}
	}
	return ctx.Serve(err)
}

func (h TransactionHistoryHandler) readDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.RequestGetDetailSO

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	if req.Session, err = service.CustomerSession(ctx); err == nil {
		ctx.ResponseData, err = h.ServicesTransactionHistory.GetDetail(ctx.Request().Context(), req)
		//ctx.Data(res)
		if err != nil {
			return ctx.Serve(err)
		}
	}
	return ctx.Serve(err)
}

func (h TransactionHistoryHandler) readDetailSI(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.RequestGetInvoiceDetail

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	if req.Session, err = service.CustomerSession(ctx); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesTransactionHistory.GetInvoiceDetail(ctx.Request().Context(), req)

	return ctx.Serve(err)
}

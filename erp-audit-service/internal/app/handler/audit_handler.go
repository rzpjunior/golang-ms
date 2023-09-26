package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type AuditHandler struct {
	Option       global.HandlerOptions
	ServiceAudit service.IAuditService
}

// URLMapping implements router.RouteHandlers
func (h *AuditHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceAudit = service.NewAuditService()

	r.GET("", h.Index)
	r.POST("", h.Create)
}

func (h AuditHandler) Index(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	auditType := ctx.GetParamString("type")
	referenceID := ctx.GetParamString("reference_id")

	var auditLogs []dto.AuditResponseGet
	var total int64
	auditLogs, total, err = h.ServiceAudit.Get(ctx.Request().Context(), page.Start, page.Limit, auditType, referenceID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(auditLogs, total, page)

	return ctx.Serve(err)
}

func (h AuditHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.AuditRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceAudit.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

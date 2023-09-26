// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"strings"
	"time"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("stc_rdl"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.PUT("/update_commited", h.updateCommited, auth.Authorized("stc_upd_cmt"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Stock
	var total int64

	if data, total, e = repository.GetStocks(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var pID int64
	productId := ctx.QueryParam("product_id")
	if strings.TrimSpace(productId) != "" {
		pID, _ = common.Decrypt(productId)
	}

	var data []*model.Stock
	var total int64

	orderChannelRestriction := strings.TrimSpace(ctx.QueryParam("order_channel_restriction"))

	if data, total, e = repository.GetFilterStocksWithProductGroup(rq, pID, orderChannelRestriction); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// updateCommited : function to update commited stock data
func (h *Handler) updateCommited(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r    updateCommitedRequest
		date time.Time
	)

	dateStr := ctx.QueryParam("dateStr")
	if dateStr == "" {
		date = time.Now()
	}

	if r.Session, e = auth.UserSession(ctx); e == nil {
		e = UpdateCommitedStock(r, date)
	}

	return ctx.Serve(e)
}

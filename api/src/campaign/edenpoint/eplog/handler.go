// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package eplog

import (
	"strings"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("ep_log_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("ep_log_rdd"))
	r.GET("/reset", h.resetCache, auth.Authorized("ep_log_rfs"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var merchantPointSummary []*model.MerchantPointSummary

	widgetPoint := new(model.WidgetPoint)

	if merchantPointSummary, _, err = repository.GetMerchantPointSummarys(rq); err != nil {
		err = echo.ErrNotFound
		return ctx.Serve(err)
	}

	for _, v := range merchantPointSummary {
		widgetPoint.TotalEarnPoint += v.EarnedPoint
		widgetPoint.TotalRedeemedPoint += v.RedeemedPoint
	}

	widgetPoint.TotalPoint, widgetPoint.LastUpdated, _ = repository.GetTotalCurrentPoint()

	responseData := &model.MerchantPointList{MerchantPointSummary: merchantPointSummary, WidgetPoint: widgetPoint}

	ctx.ResponseData = responseData

	return ctx.Serve(err)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	type responseDataStruct struct {
		Merchant  *model.Merchant           `json:"merchant"`
		PointLogs []*model.MerchantPointLog `json:"point_logs"`
	}
	var (
		merchantID      int64
		filter, exclude map[string]interface{}
		period          []string
		pointLogs       []*model.MerchantPointLog
		responseData    *responseDataStruct
		merchant        *model.Merchant
	)

	if merchantID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.QueryParam("period") != "" {
		dateParam := ctx.QueryParam("period")
		period = strings.Split(dateParam, "|")
	}

	if merchant, e = repository.ValidMerchant(merchantID); e != nil {
		e = echo.ErrNotFound
	}

	filter = map[string]interface{}{"merchant_id": merchantID}
	exclude = map[string]interface{}{"transaction_type": 9}

	if len(period) > 0 {
		period[0] += " 00:00:00"
		period[1] += " 23:59:59"
		filter["created_date__between"] = period
	}
	pointLogs, _, e = repository.CheckMerchantPointLogData(filter, exclude)

	for _, v := range pointLogs {

		if v.SalesOrder != nil {
			v.SalesOrder.Read("ID")
		}

		if v.EPCampaign != nil && v.EPCampaign.ID != 0 {
			v.EPCampaign.Read("ID")
		}
	}

	merchant.FinanceArea.Read("ID")
	merchant.BusinessType.Read("ID")

	responseData = &responseDataStruct{Merchant: merchant, PointLogs: pointLogs}

	ctx.ResponseData = responseData

	return ctx.Serve(e)
}

// resetCache : handler to reset cache
func (h *Handler) resetCache(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	key := "total_current_eden_point"
	if ctx.QueryParam("key") != "" {
		key = ctx.QueryParam("key")
	}

	e = dbredis.Redis.DeleteCacheWhereLike("*" + key + "*")

	return ctx.Serve(e)
}
